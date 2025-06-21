import numpy as np
from scipy.stats import skew

def detect_skewness(data):
    """
    Detects and reports the skewness of a given dataset.

    Args:
        data (list or np.ndarray): The dataset to analyze.

    Returns:
        tuple: A tuple containing the skewness value (float or NaN) and a string
               describing the type of skewness or an error message.
    """
    if not isinstance(data, (list, np.ndarray)):
        return float('nan'), "Error: Input data must be a list or a NumPy array."

    data_array = np.asarray(data)

    if len(data_array) < 3:
        return float('nan'), "Insufficient data to calculate skewness (requires at least 3 points)."

    try:
        skewness_value = skew(data_array)
    except Exception as e:
        return float('nan'), f"Error calculating skewness: {e}"

    # Common thresholds for interpreting skewness (e.g., from Bulmer, 1979 or similar guidelines)
    # -0.5 to 0.5: Approximately symmetric
    # -1.0 to -0.5 or 0.5 to 1.0: Moderately skewed
    # < -1.0 or > 1.0: Highly skewed
    if -0.5 <= skewness_value <= 0.5:
        skewness_type = "Approximately symmetric"
    elif skewness_value > 0.5:
        skewness_type = "Positively skewed (right-skewed)"
    else: # skewness_value < -0.5
        skewness_type = "Negatively skewed (left-skewed)"

    return skewness_value, skewness_type

def generate_sample_data(data_type="normal", size=1000):
    """
    Generates different types of sample data for demonstration.

    Args:
        data_type (str): Type of data to generate ('normal', 'positive_skew', 'negative_skew').
        size (int): Number of data points to generate.

    Returns:
        np.ndarray: Generated sample data.
    """
    if data_type == "normal":
        return np.random.normal(loc=0, scale=1, size=size)
    elif data_type == "positive_skew":
        # Exponential distribution is typically positively skewed
        return np.random.exponential(scale=1.0, size=size)
    elif data_type == "negative_skew":
        # Beta distribution can be negatively skewed (e.g., alpha=5, beta=1)
        # Or transform a positive skew to negative
        return 10 - np.random.exponential(scale=1.0, size=size)
    else:
        raise ValueError("Invalid data_type. Choose from 'normal', 'positive_skew', 'negative_skew'.")

if __name__ == "__main__":
    # Example 1: Normally distributed data (should be approximately symmetric)
    normal_data = generate_sample_data(data_type="normal", size=5000)
    skew_val_normal, skew_type_normal = detect_skewness(normal_data)
    print(f"Normal Data Skewness: {skew_val_normal:.4f} ({skew_type_normal})")

    # Example 2: Positively skewed data
    positive_skew_data = generate_sample_data(data_type="positive_skew", size=5000)
    skew_val_pos, skew_type_pos = detect_skewness(positive_skew_data)
    print(f"Positive Skew Data Skewness: {skew_val_pos:.4f} ({skew_type_pos})")

    # Example 3: Negatively skewed data
    negative_skew_data = generate_sample_data(data_type="negative_skew", size=5000)
    skew_val_neg, skew_type_neg = detect_skewness(negative_skew_data)
    print(f"Negative Skew Data Skewness: {skew_val_neg:.4f} ({skew_type_neg})")

    # Example 4: Custom list data (approximately symmetric)
    custom_data_symmetric = [1, 2, 3, 4, 5, 4, 3, 2, 1]
    skew_val_custom_sym, skew_type_custom_sym = detect_skewness(custom_data_symmetric)
    print(f"Custom Symmetric Data Skewness: {skew_val_custom_sym:.4f} ({skew_type_custom_sym})")

    # Example 5: Custom list data (positively skewed)
    custom_data_pos = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 200]
    skew_val_custom_pos, skew_type_custom_pos = detect_skewness(custom_data_pos)
    print(f"Custom Positive Skew Data Skewness: {skew_val_custom_pos:.4f} ({skew_type_custom_pos})")

    # Example 6: Custom list data (negatively skewed)
    custom_data_neg = [1, 100, 90, 80, 70, 60, 50, 40, 30, 20, 10]
    skew_val_custom_neg, skew_type_custom_neg = detect_skewness(custom_data_neg)
    print(f"Custom Negative Skew Data Skewness: {skew_val_custom_neg:.4f} ({skew_type_custom_neg})")

    # Example 7: Insufficient data
    insufficient_data = [1, 2]
    skew_val_insuf, skew_type_insuf = detect_skewness(insufficient_data)
    print(f"Insufficient Data Skewness: {skew_val_insuf} ({skew_type_insuf})")

    # Example 8: Invalid input type
    invalid_input = "this is not data"
    skew_val_invalid, skew_type_invalid = detect_skewness(invalid_input)
    print(f"Invalid Input Skewness: {skew_val_invalid} ({skew_type_invalid})")

    # Example 9: Empty list
    empty_data = []
    skew_val_empty, skew_type_empty = detect_skewness(empty_data)
    print(f"Empty Data Skewness: {skew_val_empty} ({skew_type_empty})")

# Additional implementation at 2025-06-21 04:38:21
import numpy as np
from scipy.stats import skew
import matplotlib.pyplot as plt

def detect_skewness(data, plot=True):
    """
    Detects and interprets data skewness, optionally visualizing the distribution.

    Parameters:
    data (list or np.array): The input numerical data.
    plot (bool): If True, a histogram of the data will be displayed.

    Returns:
    tuple: A tuple containing the skewness coefficient and its interpretation string.
           Returns (None, "Insufficient data") if data is too small or invalid.
    """
    if not isinstance(data, (list, np.ndarray)):
        print("Error: Input data must be a list or a NumPy array.")
        return None, "Invalid data type"

    if len(data) < 3:
        return None, "Insufficient data to calculate skewness (at least 3 points required)."

    data_array = np.array(data)

    skewness_coefficient = skew(data_array)

    if skewness_coefficient > 0.5:
        interpretation = "Positively skewed (right-skewed)"
    elif skewness_coefficient < -0.5:
        interpretation = "Negatively skewed (left-skewed)"
    else:
        interpretation = "Approximately symmetric"

    if plot:
        plt.figure(figsize=(8, 6))
        plt.hist(data_array, bins='auto', edgecolor='black', alpha=0.7)
        plt.title(f'Distribution of Data (Skewness: {skewness_coefficient:.4f})')
        plt.xlabel('Value')
        plt.ylabel('Frequency')
        plt.grid(axis='y', alpha=0.75)
        plt.show()

    return skewness_coefficient, interpretation

if __name__ == "__main__":
    data1 = [10, 20, 30, 35, 40, 45, 50, 55, 60, 70, 80, 90, 100, 200, 300, 500]
    print("--- Data Set 1 (Positively Skewed Example) ---")
    skew1, interp1 = detect_skewness(data1)
    if skew1 is not None:
        print(f"Skewness Coefficient: {skew1:.4f}")
        print(f"Interpretation: {interp1}")
    else:
        print(f"Could not calculate skewness: {interp1}")
    print("-" * 40)

    data2 = [10, 20, 30, 40, 50, 60, 70, 80, 90, 95, 98, 99, 100, 100, 100]
    print("\n--- Data Set 2 (Negatively Skewed Example) ---")
    skew2, interp2 = detect_skewness(data2)
    if skew2 is not None:
        print(f"Skewness Coefficient: {skew2:.4f}")
        print(f"Interpretation: {interp2}")
    else:
        print(f"Could not calculate skewness: {interp2}")
    print("-" * 40)

    data3 = [160, 162, 165, 168, 170, 172, 175, 178, 180, 182, 185, 188, 190]
    print("\n--- Data Set 3 (Approximately Symmetric Example) ---")
    skew3, interp3 = detect_skewness(data3)
    if skew3 is not None:
        print(f"Skewness Coefficient: {skew3:.4f}")
        print(f"Interpretation: {interp3}")
    else:
        print(f"Could not calculate skewness: {interp3}")
    print("-" * 40)

    data4 = [1, 2]
    print("\n--- Data Set 4 (Insufficient Data Example) ---")
    skew4, interp4 = detect_skewness(data4, plot=False)
    if skew4 is not None:
        print(f"Skewness Coefficient: {skew4:.4f}")
    print(f"Interpretation: {interp4}")
    print("-" * 40)

    data5 = "this is not a list of numbers"
    print("\n--- Data Set 5 (Invalid Data Type Example) ---")
    skew5, interp5 = detect_skewness(data5, plot=False)
    if skew5 is not None:
        print(f"Skewness Coefficient: {skew5:.4f}")
    print(f"Interpretation: {interp5}")
    print("-" * 40)

# Additional implementation at 2025-06-21 04:39:12
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from scipy import stats

def detect_skewness(data, plot_histogram=True, threshold=0.5, title="Data Skewness Analysis"):
    """
    Detects data skewness, calculates the skewness coefficient, and optionally
    plots a histogram.

    Args:
        data (list, np.ndarray, pd.Series): The input numerical data.
        plot_histogram (bool): If True, a histogram of the data will be plotted.
        threshold (float): The absolute skewness coefficient value above which
                           data is considered skewed.
        title (str): Title for the plot if plot_histogram is True.

    Returns:
        dict: A dictionary containing:
              - 'skewness_coefficient': The calculated skewness coefficient.
              - 'skewness_type': A string indicating the type of skewness
                                 ('Symmetric', 'Positively skewed', 'Negatively skewed',
                                 'Undefined (constant data)', 'Not enough data').
              - 'is_skewed': Boolean, True if the absolute skewness is >= threshold.
              - 'data_points': Number of valid data points analyzed.
    """
    if not isinstance(data, (list, np.ndarray, pd.Series)):
        raise TypeError("Input data must be a list, numpy array, or pandas Series.")

    # Convert to pandas Series first for robust handling of mixed types and NaNs
    # Then convert to numpy array, coercing non-numeric values to NaN
    try:
        data_series = pd.Series(data)
        data_array = pd.to_numeric(data_series, errors='coerce').values
    except Exception as e:
        raise ValueError(f"Could not convert input data to numeric array: {e}")

    # Filter out NaN values
    data_array = data_array[~np.isnan(data_array)]

    num_points = len(data_array)

    if num_points < 2:
        return {
            'skewness_coefficient': np.nan,
            'skewness_type': 'Not enough data',
            'is_skewed': False,
            'data_points': num_points
        }

    # Check if all values are the same (standard deviation will be zero, skewness undefined)
    if np.all(data_array == data_array[0]):
        return {
            'skewness_coefficient': np.nan,
            'skewness_type': 'Undefined (constant data)',
            'is_skewed': False,
            'data_points': num_points
        }

    skew_value = stats.skew(data_array)
    skewness_type = ""
    is_skewed = False

    if abs(skew_value) < threshold:
        skewness_type = 'Symmetric'
        is_skewed = False
    elif skew_value > threshold:
        skewness_type = 'Positively skewed (right-skewed)'
        is_skewed = True
    else: # skew_value < -threshold
        skewness_type = 'Negatively skewed (left-skewed)'
        is_skewed = True

    if plot_histogram:
        plt.figure(figsize=(10, 6))
        plt.hist(data_array, bins='auto', edgecolor='black', alpha=0.7)
        plt.title(f'{title}\nSkewness: {skew_value:.4f} ({skewness_type})')
        plt.xlabel('Value')
        plt.ylabel('Frequency')
        plt.grid(True, linestyle='--', alpha=0.6)
        
        # Add mean and median lines for visual inspection of skewness
        data_mean = np.mean(data_array)
        data_median = np.median(data_array)
        plt.axvline(data_mean, color='red', linestyle='dashed', linewidth=1, label=f'Mean: {data_mean:.2f}')
        plt.axvline(data_median, color='green', linestyle='dashed', linewidth=1, label=f'Median: {data_median:.2f}')
        plt.legend()
        plt.show()

    return {
        'skewness_coefficient': skew_value,
        'skewness_type': skewness_type,
        'is_skewed': is_skewed,
        'data_points': num_points
    }

if __name__ == "__main__":
    print("--- Skewness Detection Program ---")

    # Example 1: Normally distributed data (symmetric)
    data_normal = np.random.normal(loc=0, scale=1, size=1000)
    print("\nAnalyzing Normally Distributed Data:")
    result_normal = detect_skewness(data_normal, title="Normally Distributed Data")
    print(result_normal)

    # Example 2: Positively skewed data (e.g., exponential distribution)
    data_positive_skew = np.random.exponential(scale=2, size=1000)
    print("\nAnalyzing Positively Skewed Data (Exponential):")
    result_positive = detect_skewness(data_positive_skew, title="Positively Skewed Data (Exponential)")
    print(result_positive)

    # Example 3: Negatively skewed data (e.g., 10 - exponential)
    data_negative_skew = 10 - np.random.exponential(scale=2, size=1000)
    print("\nAnalyzing Negatively Skewed Data (10 - Exponential):")
    result_negative = detect_skewness(data_negative_skew, title="Negatively Skewed Data (10 - Exponential)")
    print(result_negative)

    # Example 4: Data with NaNs and mixed types (list input)
    data_with_nan = [1, 2, 3, 4, np.nan, 5, 6, 'invalid', 7, 8, 9, 10, None, 11, 12, 13, 14, 15]
    print("\nAnalyzing Data with NaNs, None, and Non-numeric (List Input):")
    result_nan = detect_skewness(data_with_nan, title="Data with NaNs, None, and Non-numeric")
    print(result_nan)

    # Example 5: Pandas Series input (Log-Normal distribution, typically positively skewed)
    data_pandas_series = pd.Series(np.random.lognormal(mean=0, sigma=1, size=500))
    print("\nAnalyzing Pandas Series (Log-Normal):")
    result_pandas = detect_skewness(data_pandas_series, title="Pandas Series (Log-Normal)")
    print(result_pandas)

    # Example 6: Constant data
    data_constant = [5, 5, 5, 5, 5]
    print("\nAnalyzing Constant Data:")
    result_constant = detect_skewness(data_constant, title="Constant Data")
    print(result_constant)

    # Example 7: Not enough data (single point)
    data_single_point = [10]
    print("\nAnalyzing Single Data Point:")
    result_single = detect_skewness(data_single_point, title="Single Data Point")
    print(result_single)

    # Example 8: Empty data
    data_empty = []
    print("\nAnalyzing Empty Data:")
    result_empty = detect_skewness(data_empty, title="Empty Data")
    print(result_empty)

    # Example 9: Custom threshold (e.g., a stricter threshold)
    data_custom_threshold = np.random.normal(loc=0, scale=1, size=100) # Should be symmetric
    print("\nAnalyzing Data with Custom Threshold (0.1):")
    result_custom_threshold = detect_skewness(data_custom_threshold, threshold=0.1, title="Data with Custom Threshold")
    print(result_custom_threshold)

    print("\n--- End of Program ---")