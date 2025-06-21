import os
from PIL import Image

def get_image_signature(image_path, size=(16, 16)):
    try:
        with Image.open(image_path) as img:
            img = img.convert("L")
            img = img.resize(size, Image.Resampling.LANCZOS)
            return img.tobytes()
    except Exception:
        return None

def find_duplicate_images(directory):
    image_signatures = {}
    
    for root, _, files in os.walk(directory):
        for filename in files:
            file_path = os.path.join(root, filename)
            
            if not any(file_path.lower().endswith(ext) for ext in ['.png', '.jpg', '.jpeg', '.gif', '.bmp', '.tiff', '.webp']):
                continue

            signature = get_image_signature(file_path)
            
            if signature is not None:
                if signature not in image_signatures:
                    image_signatures[signature] = []
                image_signatures[signature].append(file_path)
    
    duplicate_groups = []
    for signature, paths in image_signatures.items():
        if len(paths) > 1:
            duplicate_groups.append(paths)
            
    return duplicate_groups

# Additional implementation at 2025-06-21 00:17:52
import os
from PIL import Image, ImageChops
import imagehash
import math

class ImageDuplicateFinder:
    def __init__(self, directory, resize_dim=(256, 256), pixel_comparison_threshold=1000, hash_size=8, hash_threshold=5):
        """
        Initializes the ImageDuplicateFinder.

        Args:
            directory (str): The path to the directory to scan for images.
            resize_dim (tuple): Dimensions (width, height) to resize images for pixel comparison.
                                 Set to None to disable resizing for pixel comparison.
            pixel_comparison_threshold (int): Max MSE difference for images to be considered duplicates
                                              using pixel comparison. Lower is stricter.
            hash_size (int): The size of the hash (e.g., 8 for 8x8 hash).
            hash_threshold (int): The maximum Hamming distance for hashes to be considered duplicates.
                                  Lower is stricter.
        """
        if not os.path.isdir(directory):
            raise ValueError(f"Directory not found: {directory}")
        self.directory = directory
        self.resize_dim = resize_dim
        self.pixel_comparison_threshold = pixel_comparison_threshold
        self.hash_size = hash_size
        self.hash_threshold = hash_threshold
        self.image_extensions = ('.png', '.jpg', '.jpeg', '.gif', '.bmp', '.tiff')

    def _load_and_preprocess_image(self, image_path, for_pixel_comparison=True):
        """
        Loads an image, converts it to RGB, and optionally resizes it.
        """
        try:
            img = Image.open(image_path).convert("RGB")
            if for_pixel_comparison and self.resize_dim:
                img = img.resize(self.resize_dim, Image.Resampling.LANCZOS)
            return img
        except Exception:
            return None

    def _compare_pixels_mse(self, img1, img2):
        """
        Compares two PIL images pixel by pixel using Mean Squared Error (MSE).
        Assumes images are already preprocessed (same size, same mode 'RGB').
        """
        if img1.size != img2.size or img1.mode != img2.mode:
            return float('inf')

        if img1 == img2:
            return 0

        pixels1 = img1.load()
        pixels2 = img2.load()
        
        width, height = img1.size
        sq_diff_sum = 0

        for y in range(height):
            for x in range(width):
                r1, g1, b1 = pixels1[x, y]
                r2, g2, b2 = pixels2[x, y]
                sq_diff_sum += (r1 - r2)**2 + (g1 - g2)**2 + (b1 - b2)**2

        mse = sq_diff_sum / (width * height * 3)
        return mse

    def _get_image_hash(self, img_path, hash_type='phash'):
        """
        Generates a perceptual hash for an image.
        """
        img = self._load_and_preprocess_image(img_path, for_pixel_comparison=False)
        if img is None:
            return None

        try:
            if hash_type == 'ahash':
                return imagehash.average_hash(img, hash_size=self.hash_size)
            elif hash_type == 'phash':
                return imagehash.phash(img, hash_size=self.hash_size)
            elif hash_type == 'dhash':
                return imagehash.dhash(img, hash_size=self.hash_size)
            elif hash_type == 'whash':
                return imagehash.whash(img, hash_size=self.hash_size)
            else:
                return imagehash.phash(img, hash_size=self.hash_size)
        except Exception:
            return None

    def find_duplicates_pixel_comparison(self):
        """
        Finds duplicate images in the directory using pixel-by-pixel comparison (MSE).
        Returns a dictionary where keys are original image paths and values are lists of duplicate paths.
        """
        image_paths = [os.path.join(self.directory, f) for f in os.listdir(self.directory)
                       if f.lower().endswith(self.image_extensions) and os.path.isfile(os.path.join(self.directory, f))]

        processed_images = {}
        for path in image_paths:
            img = self._load_and_preprocess_image(path, for_pixel_comparison=True)
            if img:
                processed_images[path] = img

        duplicates = {}
        compared_pairs = set()

        paths_list = list(processed_images.keys())
        for i in range(len(paths_list)):
            path1 = paths_list[i]
            img1 = processed_images[path1]
            if img1 is None:
                continue
            
            current_duplicates_for_path1 = []
            for j in range(i + 1, len(paths_list)):
                path2 = paths_list[j]
                img2 = processed_images[path2]
                if img2 is None:
                    continue

                pair = tuple(sorted((path1, path2)))
                if pair in compared_pairs:
                    continue
                compared_pairs.add(pair)

                mse = self._compare_pixels_mse(img1, img2)
                if mse <= self.pixel_comparison_threshold:
                    current_duplicates_for_path1.append(path2)
            
            if current_duplicates_for_path1:
                duplicates[path1] = current_duplicates_for_path1
        
        return duplicates

    def find_duplicates_perceptual_hash(self, hash_type='phash'):
        """
        Finds duplicate images in the directory using perceptual hashing.
        Returns a list of lists, where each inner list contains paths of images
        that are considered duplicates of each other.
        """
        image_paths = [os.path.join(self.directory, f) for f in os.listdir(self.directory)
                       if f.lower().endswith(self.image_extensions) and os.path.isfile(os.path.join(self.directory, f))]

        hashes = {}
        for path in image_paths:
            img_hash = self._get_image_hash(path, hash_type)
            if img_hash:
                if img_hash not in hashes:
                    hashes[img_hash] = []
                hashes[img_hash].append(path)

        duplicate_groups = []
        processed_hashes = set()

        hash_list = list(hashes.keys())
        for i in range(len(hash_list)):
            h1 = hash_list[i]
            if h1 in processed_hashes:
                continue

            current_group = set(hashes[h1])
            processed_hashes.add(h1)

            for j in range(i + 1, len(hash_list)):
                h2 = hash_list[j]
                if h2 in processed_hashes:
                    continue

                distance = h1 - h2
                if distance <= self.hash_threshold:
                    current_group.update(hashes[h2])
                    processed_hashes.add(h2)
            
            if len(current_group) > 1:
                duplicate_groups.append(list(current_group))
        
        return duplicate_groups

    def find_duplicates(self, method='perceptual_hash', hash_type='phash'):
        """
        Main method to find duplicates using the specified method.

        Args:
            method (str): 'pixel_comparison' or 'perceptual_hash'.
            hash_type (str): 'ahash', 'phash', 'dhash', 'whash' (only for 'perceptual_hash' method).

        Returns:
            dict or list:
                If 'pixel_comparison': A dictionary where keys are original image paths and values are lists of duplicate paths.
                If 'perceptual_hash': A list of lists, where each inner list contains paths of images
                                     that are considered duplicates of each other.
        """
        if method == 'pixel_comparison':
            return self.find_duplicates_pixel_comparison()
        elif method == 'perceptual_hash':
            return self.find_duplicates_perceptual_hash(hash_type)
        else:
            raise ValueError("Invalid method. Choose 'pixel_comparison' or 'perceptual_hash'.")