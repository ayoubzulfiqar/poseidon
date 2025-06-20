def calculate_percentile(data, percentile):
    n = len(data)
    if n == 0:
        return None
    k = (n - 1) * (percentile / 100.0)
    f = int(k)
    c = k - f
    if f + 1 < n:
        return data[f] + (data[f+1] - data[f]) * c
    else:
        return data[f]

def find_outliers(data):
    if not data or len(data) < 2:
        return []

    sorted_data = sorted(data)
    
    q1 = calculate_percentile(sorted_data, 25)
    q3 = calculate_percentile(sorted_data, 75)
    
    iqr = q3 - q1
    
    lower_bound = q1 - 1.5 * iqr
    upper_bound = q3 + 1.5 * iqr
    
    outliers = [x for x in data if x < lower_bound or x > upper_bound]
    
    return outliers

if __name__ == "__main__":
    data_set1 = [1, 2, 3, 4, 5, 6, 7, 8, 9, 100]
    data_set2 = [10, 12, 12, 13, 12, 11, 14, 13, 15, 10, 6, 13, 12, 100, 1, 200]
    data_set3 = [1, 2, 3, 4, 5]
    data_set4 = []
    data_set5 = [10]
    data_set6 = [1, 1, 1, 100, 1, 1, 1]

    print(f"Data Set: {data_set1}")
    print(f"Outliers: {find_outliers(data_set1)}")

    print(f"Data Set: {data_set2}")
    print(f"Outliers: {find_outliers(data_set2)}")

    print(f"Data Set: {data_set3}")
    print(f"Outliers: {find_outliers(data_set3)}")
    
    print(f"Data Set: {data_set4}")
    print(f"Outliers: {find_outliers(data_set4)}")

    print(f"Data Set: {data_set5}")
    print(f"Outliers: {find_outliers(data_set5)}")

    print(f"Data Set: {data_set6}")
    print(f"Outliers: {find_outliers(data_set6)}")

# Additional implementation at 2025-06-19 23:38:44
import numpy as np
from scipy import stats
import matplotlib.pyplot as plt

def find_outliers_zscore(data, threshold=3.0):
    """
    Finds outliers in a dataset using the Z-score method.

    Args:
        data (np.array or list): The input data.
        threshold (float): The Z-score threshold for identifying outliers.

    Returns:
        tuple: A tuple containing:
            - np.array: Indices of the identified outliers.
            - np.array: Values of the identified outliers.
            - np.array: Boolean mask indicating outliers.
    """
    data = np.asarray(data)
    z_scores = np.abs(stats.zscore(data))
    outlier_mask = z_scores > threshold
    outlier_indices = np.where(outlier_mask)[0]
    outlier_values = data[outlier_indices]
    return outlier_indices, outlier_values, outlier_mask

def find_outliers_iqr(data, k=1.5):
    """
    Finds outliers in a dataset using the Interquartile Range (IQR) method.

    Args:
        data (np.array or list): The input data.
        k (float): The multiplier for the IQR to define the outlier bounds.
                   Commonly 1.5 for mild outliers, 3.0 for extreme outliers.

    Returns:
        tuple: A tuple containing:
            - np.array: Indices of the identified outliers.
            - np.array: Values of the identified outliers.
            - np.array: Boolean mask indicating outliers.
    """
    data = np.asarray(data)
    Q1 = np.percentile(data, 25)
    Q3 = np.percentile(data, 75)
    IQR = Q3 - Q1
    lower_bound = Q1 - k * IQR
    upper_bound = Q3 + k * IQR
    outlier_mask = (data < lower_bound) | (data > upper_bound)
    outlier_indices = np.where(outlier_mask)[0]
    outlier_values = data[outlier_indices]
    return outlier_indices, outlier_values, outlier_mask

def plot_outliers(data, outlier_mask, method_name):
    """
    Plots the data and highlights the identified outliers.

    Args:
        data (np.array): The input data.
        outlier_mask (np.array): Boolean mask indicating outliers.
        method_name (str): Name of the outlier detection method for the plot title.
    """
    plt.figure(figsize=(10, 6))
    
    plt.scatter(np.arange(len(data)), data, color='blue', label='Normal Data Points', s=50)
    
    outlier_indices = np.where(outlier_mask)[0]
    outlier_values = data[outlier_indices]
    plt.scatter(outlier_indices, outlier_values, color='red', label='Outliers', s=100, edgecolors='black', zorder=5)
    
    plt.title(f'Outlier Detection using {method_name} Method')
    plt.xlabel('Data Index')
    plt.ylabel('Value')
    plt.legend()
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.show()

if __name__ == "__main__":
    np.random.seed(42)
    data = np.random.normal(loc=50, scale=10, size=100)
    data = np.append(data, [1, 2, 90, 95, 100, -10, -5])
    data = np.sort(data)

    z_outlier_indices, z_outlier_values, z_outlier_mask = find_outliers_zscore(data, threshold=2.5)
    print(z_outlier_values)
    plot_outliers(data, z_outlier_mask, "Z-score")

    iqr_outlier_indices, iqr_outlier_values, iqr_outlier_mask = find_outliers_iqr(data, k=1.5)
    print(iqr_outlier_values)
    plot_outliers(data, iqr_outlier_mask, "IQR")

# Additional implementation at 2025-06-19 23:39:45
import numpy as np

class OutlierDetector:
    """
    A class to detect statistical outliers in a dataset using various methods.
    """

    def __init__(self, data):
        """
        Initializes the OutlierDetector with a dataset.

        Args:
            data (list or numpy.ndarray): The numerical dataset to analyze.
        """
        if not isinstance(data, (list, np.ndarray)):
            raise TypeError("Input data must be a list or numpy array.")
        if not data:
            self.data = np.array([])
        else:
            self.data = np.array(data, dtype=float)
            if not np.issubdtype(self.data.dtype, np.number):
                raise ValueError("Input data must contain only numerical values.")

    def _check_data_sufficiency(self, min_len=2):
        """Helper to check if data is sufficient for calculations."""
        return len(self.data) >= min_len

    def find_outliers_zscore(self, threshold=3.0):
        """
        Detects outliers using the Z-score method.
        Assumes data is normally distributed.

        Args:
            threshold (float): The Z-score threshold. Values with an absolute
                               Z-score greater than this threshold are considered outliers.

        Returns:
            tuple: A tuple containing:
                   - list: Outlier values.
                   - list: Indices of outlier values in the original data.
        """
        if not self._check_data_sufficiency(min_len=2):
            return [], []

        mean = np.mean(self.data)
        std_dev = np.std(self.data)

        if std_dev == 0:
            return [], []

        z_scores = np.abs((self.data - mean) / std_dev)
        outlier_indices = np.where(z_scores > threshold)[0]
        outlier_values = self.data[outlier_indices].tolist()

        return outlier_values, outlier_indices.tolist()

    def find_outliers_iqr(self):
        """
        Detects outliers using the Interquartile Range (IQR) method.
        This method is robust to non-normal distributions.

        Returns:
            tuple: A tuple containing:
                   - list: Outlier values.
                   - list: Indices of outlier values in the original data.
        """
        if not self._check_data_sufficiency(min_len=2):
            return [], []

        Q1 = np.percentile(self.data, 25)
        Q3 = np.percentile(self.data, 75)
        IQR = Q3 - Q1

        if IQR == 0:
            return [], []

        lower_bound = Q1 - 1.5 * IQR
        upper_bound = Q3 + 1.5 * IQR

        outlier_indices = np.where((self.data < lower_bound) | (self.data > upper_bound))[0]
        outlier_values = self.data[outlier_indices].tolist()

        return outlier_values, outlier_indices.tolist()

    def get_data_summary(self):
        """
        Provides a basic statistical summary of the dataset.

        Returns:
            dict: A dictionary containing summary statistics.
        """
        if not self._check_data_sufficiency(min_len=1):
            return {
                "count": 0,
                "mean": None,
                "median": None,
                "std_dev": None,
                "min": None,
                "max": None,
                "Q1": None,
                "Q3": None,
                "IQR": None
            }

        summary = {
            "count": len(self.data),
            "mean": np.mean(self.data),
            "median": np.median(self.data),
            "std_dev": np.std(self.data),
            "min": np.min(self.data),
            "max": np.max(self.data),
            "Q1": np.percentile(self.data, 25),
            "Q3": np.percentile(self.data, 75),
        }
        summary["IQR"] = summary["Q3"] - summary["Q1"]
        return summary

    def remove_outliers(self, method='iqr', threshold=3.0):
        """
        Returns a new dataset with detected outliers removed.

        Args:
            method (str): The outlier detection method to use ('zscore' or 'iqr').
                          Defaults to 'iqr'.
            threshold (float): Z-score threshold if 'zscore' method is chosen.

        Returns:
            numpy.ndarray: A new array with outliers removed.
        """
        if method == 'zscore':
            _, outlier_indices = self.find_outliers_zscore(threshold)
        elif method == 'iqr':
            _, outlier_indices = self.find_outliers_iqr()
        else:
            raise ValueError("Invalid method. Choose 'zscore' or 'iqr'.")

        non_outlier_mask = np.ones(len(self.data), dtype=bool)
        non_outlier_mask[outlier_indices] = False

        return self.data[non_outlier_mask]