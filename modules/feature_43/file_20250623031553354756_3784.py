import math

def calculate_mean(data):
    """Calculates the mean of a list of numbers."""
    if not data:
        return 0.0
    return sum(data) / len(data)

def calculate_standard_deviation(data, mean_val=None):
    """Calculates the standard deviation of a list of numbers."""
    if len(data) < 2:
        return 0.0
    if mean_val is None:
        mean_val = calculate_mean(data)
    variance = sum([(x - mean_val) ** 2 for x in data]) / (len(data) - 1)
    return math.sqrt(variance)

def calculate_covariance(data1, data2, mean1=None, mean2=None):
    """Calculates the covariance between two lists of numbers."""
    if len(data1) != len(data2) or not data1:
        return 0.0
    if mean1 is None:
        mean1 = calculate_mean(data1)
    if mean2 is None:
        mean2 = calculate_mean(data2)
    
    covariance_sum = sum([(data1[i] - mean1) * (data2[i] - mean2) for i in range(len(data1))])
    return covariance_sum / (len(data1) - 1)

def calculate_correlation_coefficient(data1, data2):
    """Calculates the Pearson correlation coefficient between two lists of numbers."""
    if len(data1) != len(data2) or len(data1) < 2:
        return 0.0 # Or raise an error for insufficient data

    mean1 = calculate_mean(data1)
    mean2 = calculate_mean(data2)

    std_dev1 = calculate_standard_deviation(data1, mean1)
    std_dev2 = calculate_standard_deviation(data2, mean2)

    if std_dev1 == 0 or std_dev2 == 0:
        return 0.0 # Or handle division by zero as appropriate (e.g., return NaN or 0)

    covariance = calculate_covariance(data1, data2, mean1, mean2)

    return covariance / (std_dev1 * std_dev2)

if __name__ == '__main__':
    # Example Usage:
    # Data for variable X
    X = [10, 20, 30, 40, 50]
    # Data for variable Y
    Y = [15, 25, 35, 45, 55]

    # Calculate correlation coefficient
    correlation = calculate_correlation_coefficient(X, Y)
    print(f"Correlation coefficient between X and Y: {correlation}")

    # Another example with a negative correlation
    A = [1, 2, 3, 4, 5]
    B = [5, 4, 3, 2, 1]
    correlation_ab = calculate_correlation_coefficient(A, B)
    print(f"Correlation coefficient between A and B: {correlation_ab}")

    # Example with no correlation (or very low)
    C = [1, 2, 3, 4, 5]
    D = [1, 1, 5, 5, 1]
    correlation_cd = calculate_correlation_coefficient(C, D)
    print(f"Correlation coefficient between C and D: {correlation_cd}")

    # Example with identical data (perfect positive correlation)
    E = [1, 2, 3, 4, 5]
    F = [1, 2, 3, 4, 5]
    correlation_ef = calculate_correlation_coefficient(E, F)
    print(f"Correlation coefficient between E and F: {correlation_ef}")

    # Example with constant data (std dev will be 0)
    G = [10, 10, 10, 10, 10]
    H = [1, 2, 3, 4, 5]
    correlation_gh = calculate_correlation_coefficient(G, H)
    print(f"Correlation coefficient between G and H: {correlation_gh}")

    # Example with empty lists
    I = []
    J = []
    correlation_ij = calculate_correlation_coefficient(I, J)
    print(f"Correlation coefficient between I and J: {correlation_ij}")

    # Example with different lengths
    K = [1, 2, 3]
    L = [4, 5]
    correlation_kl = calculate_correlation_coefficient(K, L)
    print(f"Correlation coefficient between K and L: {correlation_kl}")

# Additional implementation at 2025-06-23 03:16:29
import numpy as np
import pandas as pd

def calculate_pearson_correlation(data1, data2):
    """
    Calculates the Pearson correlation coefficient between two 1D arrays.
    Handles NaN values by removing pairs where either value is NaN.
    Returns np.nan if there are fewer than 2 valid data points.
    """
    s1 = pd.Series(data1)
    s2 = pd.Series(data2)

    # Combine into a DataFrame to drop rows with NaNs in either series
    temp_df = pd.DataFrame({'s1': s1, 's2': s2}).dropna()

    if len(temp_df) < 2:
        # Not enough data points to calculate correlation after dropping NaNs
        return np.nan

    # Calculate correlation using pandas built-in method on the cleaned data
    return temp_df['s1'].corr(temp_df['s2'])

def calculate_correlation_matrix(dataframe):
    """
    Calculates the Pearson correlation matrix for all numerical columns in a DataFrame.
    Handles NaN values by default (pairwise deletion for each correlation calculation).
    Non-numeric columns are automatically excluded.
    """
    return dataframe.corr(method='pearson')

if __name__ == "__main__":
    # --- Example Data ---
    # Sample data with some missing values (NaN) and a non-numeric column
    data = {
        'Variable_A': [10, 20, 30, 40, 50, 60, 70, 80, 90, 100],
        'Variable_B': [12, 22, 35, 41, 55, 63, 78, 81, 92, 105],
        'Variable_C': [5, 10, 15, 20, 25, np.nan, 35, 40, 45, 50], # Missing value
        'Variable_D': [100, 90, 80, 70, 60, 50, 40, 30, 20, np.nan], # Missing value
        'Variable_E': [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
        'Non_Numeric': ['X', 'Y', 'Z', 'A', 'B', 'C', 'D', 'E', 'F', 'G'] # Non-numeric column
    }
    df = pd.DataFrame(data)

    print("--- Original DataFrame ---")
    print(df)
    print("\n" + "="*40 + "\n")

    # --- Additional Functionality 1: Calculate Pearson correlation between two specific variables ---
    # This function demonstrates handling of NaNs for pairwise correlation.
    print("--- Pearson Correlation between two specific variables (Variable_A and Variable_B) ---")
    corr_ab = calculate_pearson_correlation(df['Variable_A'], df['Variable_B'])
    print(f"Correlation (A, B): {corr_ab:.4f}")

    print("\n--- Pearson Correlation between Variable_A and Variable_C (with NaN in C) ---")
    corr_ac = calculate_pearson_correlation(df['Variable_A'], df['Variable_C'])
    print(f"Correlation (A, C): {corr_ac:.4f}")

    print("\n--- Pearson Correlation between Variable_C and Variable_D (with NaNs in both) ---")
    corr_cd = calculate_pearson_correlation(df['Variable_C'], df['Variable_D'])
    print(f"Correlation (C, D): {corr_cd:.4f}")

    # Example with very few non-NaN pairs, resulting in NaN correlation
    print("\n--- Pearson Correlation with very few valid pairs (should be NaN) ---")
    sparse_data1 = [1, 2, np.nan, np.nan, np.nan, np.nan, np.nan, np.nan, np.nan, 10]
    sparse_data2 = [np.nan, np.nan, 3, 4, np.nan, np.nan, np.nan, np.nan, 9, np.nan]
    corr_sparse = calculate_pearson_correlation(sparse_data1, sparse_data2)
    print(f"Correlation (Sparse1, Sparse2): {corr_sparse}") # Will print nan
    print("\n" + "="*40 + "\n")

    # --- Additional Functionality 2: Calculate Correlation Matrix for the entire DataFrame ---
    # This function automatically handles multiple variables, NaNs, and non-numeric columns.
    print("--- Correlation Matrix for the DataFrame ---")
    correlation_matrix = calculate_correlation_matrix(df)
    print(correlation_matrix)
    print("\nNote: Non-numeric columns ('Non_Numeric') are automatically excluded.")
    print("Note: NaN values are handled by pairwise deletion for each correlation coefficient calculation.")

# Additional implementation at 2025-06-23 03:17:18
import math

def _is_numeric(value):
    """Helper to check if a value is numeric (int or float) and not NaN."""
    return isinstance(value, (int, float)) and not (isinstance(value, float) and math.isnan(value))

def _calculate_mean(data):
    """Calculates the mean of a list of numbers."""
    if not data:
        return 0.0
    return sum(data) / len(data)

def _calculate_stdev(data, mean):
    """Calculates the sample standard deviation of a list of numbers."""
    if len(data) < 2:
        return 0.0
    variance = sum([(x - mean) ** 2 for x in data]) / (len(data) - 1)
    return math.sqrt(variance)

def _calculate_covariance(x, y, mean_x, mean_y):
    """Calculates the sample covariance between two lists of numbers."""
    if len(x) != len(y):
        raise ValueError("Input lists for covariance must have the same length.")
    if len(x) < 2:
        return 0.0 # Covariance is undefined or 0 for less than 2 points

    covariance_sum = sum([(xi - mean_x) * (yi - mean_y) for xi, yi in zip(x, y)])
    return covariance_sum / (len(x) - 1)

def pearson_correlation(x_data, y_data):
    """
    Calculates the Pearson correlation coefficient between two lists of numbers.
    This function extends basic correlation calculation by:
    - Handling non-numeric values (including float('nan')) by skipping pairs.
    - Returning float('nan') if correlation cannot be calculated (e.g., insufficient data, no variance).
    - Performing input validation for list lengths.
    """
    if len(x_data) != len(y_data):
        raise ValueError("Input lists must have the same length.")
    if not x_data:
        return float('nan') # Cannot calculate correlation for empty data

    # Filter out non-numeric values or NaNs from pairs
    filtered_x = []
    filtered_y = []
    for val_x, val_y in zip(x_data, y_data):
        if _is_numeric(val_x) and _is_numeric(val_y):
            filtered_x.append(val_x)
            filtered_y.append(val_y)

    if len(filtered_x) < 2:
        # Need at least 2 data points to calculate correlation
        return float('nan')

    mean_x = _calculate_mean(filtered_x)
    mean_y = _calculate_mean(filtered_y)

    stdev_x = _calculate_stdev(filtered_x, mean_x)
    stdev_y = _calculate_stdev(filtered_y, mean_y)

    # If either standard deviation is zero, correlation is undefined
    if stdev_x == 0 or stdev_y == 0:
        return float('nan')

    covariance = _calculate_covariance(filtered_x, filtered_y, mean_x, mean_y)

    correlation = covariance / (stdev_x * stdev_y)
    
    # Ensure correlation is within [-1, 1] due to potential floating point inaccuracies
    return max(-1.0, min(1.0, correlation))

def main():
    """
    Demonstrates the usage of the pearson_correlation function with various examples,
    showcasing its error handling and NaN/non-numeric value filtering.
    """
    print("--- Pearson Correlation Coefficient Calculator ---")

    # Example 1: Basic positive correlation
    data1_x = [1, 2, 3, 4, 5]
    data1_y = [2, 4, 5, 4, 5]
    corr1 = pearson_correlation(data1_x, data1_y)
    print(f"\nData X: {data1_x}")
    print(f"Data Y: {data1_y}")
    print(f"Correlation: {corr1:.4f}")

    # Example 2: Perfect positive correlation
    data2_x = [10, 20, 30, 40, 50]
    data2_y = [100, 200, 300, 400, 500]
    corr2 = pearson_correlation(data2_x, data2_y)
    print(f"\nData X: {data2_x}")
    print(f"Data Y: {data2_y}")
    print(f"Correlation: {corr2:.4f}")

    # Example 3: Perfect negative correlation
    data3_x = [1, 2, 3, 4, 5]
    data3_y = [5, 4, 3, 2, 1]
    corr3 = pearson_correlation(data3_x, data3_y)
    print(f"\nData X: {data3_x}")
    print(f"Data Y: {data3_y}")
    print(f"Correlation: {corr3:.4f}")

    # Example 4: No correlation (one variable has no variance)
    data4_x = [1, 2, 3, 4, 5]
    data4_y = [3, 3, 3, 3, 3]
    corr4 = pearson_correlation(data4_x, data4_y)
    print(f"\nData X: {data4_x}")
    print(f"Data Y: {data4_y}")
    print(f"Correlation: {corr4}") # Expected: nan

    # Example 5: Data with missing values (float('nan'))
    data5_x = [1, 2, float('nan'), 4, 5]
    data5_y = [2, 4, 5, float('nan'), 5]
    corr5 = pearson_correlation(data5_x, data5_y)
    print(f"\nData X: {data5_x}")
    print(f"Data Y: {data5_y}")
    print(f"Correlation: {corr5:.4f}") # Calculates on (1,2), (2,4), (5,5)

    # Example 6: Data with None values and non-numeric strings
    data6_x = [1, 2, None, 4, 'five']
    data6_y = [2, 4, 5, 4, 5]
    corr6 = pearson_correlation(data6_x, data6_y)
    print(f"\nData X: {data6_x}")
    print(f"Data Y: {data6_y}")
    print(f"Correlation: {corr6:.4f}") # Calculates on (1,2), (2,4), (4,4)

    # Example 7: Different lengths (should raise ValueError)
    data7_x = [1, 2, 3]
    data7_y = [1, 2]
    try:
        corr7 = pearson_correlation(data7_x, data7_y)
        print(f"\nData X: {data7_x}")
        print(f"Data Y: {data7_y}")
        print(f"Correlation: {corr7}")
    except ValueError as e:
        print(f"\nError for different lengths: {e}")

    # Example 8: Empty lists
    data8_x = []
    data8_y = []
    corr8 = pearson_correlation(data8_x, data8_y)
    print(f"\nData X: {data8_x}")
    print(f"Data Y: {data8_y}")
    print(f"Correlation: {corr8}") # Expected: nan

    # Example 9: Single data point (should return nan as correlation requires at least 2 points)
    data9_x = [1]
    data9_y = [2]
    corr9 = pearson_correlation(data9_x, data9_y)
    print(f"\nData X: {data9_x}")
    print(f"Data Y: {data9_y}")
    print(f"Correlation: {corr9}") # Expected: nan

if __name__ == "__main__":
    main()

# Additional implementation at 2025-06-23 03:18:08
import math

def calculate_mean(data):
    """Calculates the arithmetic mean of a list of numbers."""
    if not data:
        return 0.0
    return sum(data) / len(data)

def pearson_correlation(data1, data2):
    """
    Calculates the Pearson correlation coefficient between two lists of numbers.

    Args:
        data1 (list): The first list of numerical data.
        data2 (list): The second list of numerical data.

    Returns:
        float: The Pearson correlation coefficient.

    Raises:
        ValueError: If the input lists are not of the same length,
                    if there are fewer than two data points, or
                    if one or both datasets have zero variance.
    """
    if len(data1) != len(data2):
        raise ValueError("Input lists must have the same length.")
    if len(data1) < 2:
        raise ValueError("At least two data points are required to calculate correlation.")

    mean1 = calculate_mean(data1)
    mean2 = calculate_mean(data2)

    # Calculate the numerator (sum of products of deviations)
    numerator = sum([(x - mean1) * (y - mean2) for x, y in zip(data1, data2)])

    # Calculate the sum of squared deviations for each dataset
    sum_sq_dev1 = sum([(x - mean1) ** 2 for x in data1])
    sum_sq_dev2 = sum([(y - mean2) ** 2 for y in data2])

    # Calculate the denominator (product of square roots of sum of squared deviations)
    denominator = math.sqrt(sum_sq_dev1) * math.sqrt(sum_sq_dev2)

    if denominator == 0:
        # This occurs if one or both datasets have zero variance (all values are identical).
        # In such cases, the correlation is undefined.
        raise ValueError("Correlation is undefined due to zero variance in one or both datasets.")

    return numerator / denominator

def get_numeric_input(prompt):
    """
    Prompts the user for numerical input and converts it to a list of floats.
    Handles invalid input gracefully.
    """
    while True:
        user_input = input(prompt)
        try:
            # Replace commas with spaces to allow both comma-separated and space-separated input
            numbers_str = user_input.replace(',', ' ').split()
            numbers = [float(n) for n in numbers_str]
            if not numbers:
                print("No numbers entered. Please try again.")
                continue
            return numbers
        except ValueError:
            print("Invalid input. Please enter numbers separated by spaces or commas.")

def main():
    """
    Main function to run the Pearson Correlation Coefficient Calculator.
    Provides user interaction, input validation, and result interpretation.
    """
    print("Pearson Correlation Coefficient Calculator")
    print("----------------------------------------")

    while True:
        try:
            data_x = get_numeric_input("Enter the first set of numbers (e.g., 1 2 3 4 5 or 1,2,3,4,5): ")
            data_y = get_numeric_input("Enter the second set of numbers (e.g., 6 7 8 9 10 or 6,7,8,9,10): ")

            correlation = pearson_correlation(data_x, data_y)
            print(f"\nCorrelation Coefficient (r): {correlation:.4f}")

            # Additional functionality: Interpretation of the correlation coefficient
            if correlation >= 0.7:
                print("Interpretation: Strong positive linear correlation.")
            elif correlation >= 0.3:
                print("Interpretation: Moderate positive linear correlation.")
            elif correlation > 0:
                print("Interpretation: Weak positive linear correlation.")
            elif correlation == 0:
                print("Interpretation: No linear correlation.")
            elif correlation <= -0.7:
                print("Interpretation: Strong negative linear correlation.")
            elif correlation <= -0.3:
                print("Interpretation: Moderate negative linear correlation.")
            else: # correlation < 0
                print("Interpretation: Weak negative linear correlation.")

        except ValueError as e:
            print(f"Error: {e}")
        except Exception as e:
            print(f"An unexpected error occurred: {e}")

        # Additional functionality: Allow multiple calculations
        again = input("\nDo you want to calculate another correlation? (yes/no): ").lower()
        if again != 'yes':
            break

if __name__ == '__main__':
    main()