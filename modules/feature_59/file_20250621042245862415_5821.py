def evaluate_polynomial(coefficients, x_value):
    """
    Evaluates a polynomial for a given x_value.

    The polynomial is represented by a list of coefficients, where
    coefficients[i] is the coefficient of x^i.
    For example, [a0, a1, a2] represents the polynomial a2*x^2 + a1*x + a0.
    """
    result = 0
    for i, coeff in enumerate(coefficients):
        result += coeff * (x_value ** i)
    return result

# Additional implementation at 2025-06-21 04:23:38
import collections

class Polynomial:
    def __init__(self, coefficients):
        """
        Initializes a Polynomial object.
        Coefficients can be provided as:
        - A list: [c0, c1, c2, ...] where ci is the coefficient of x^i.
        - A dictionary: {power: coefficient, ...}.
        """
        self.coeffs = collections.defaultdict(float)
        if isinstance(coefficients, list):
            for i, coeff in enumerate(coefficients):
                if coeff != 0:
                    self.coeffs[i] = float(coeff)
        elif isinstance(coefficients, dict):
            for power, coeff in coefficients.items():
                if coeff != 0:
                    self.coeffs[int(power)] = float(coeff)
        else:
            raise TypeError("Coefficients must be a list or a dictionary.")
        
        self._clean_coefficients()

    def _clean_coefficients(self):
        """Removes zero coefficients and ensures consistency."""
        keys_to_remove = [p for p, c in self.coeffs.items() if c == 0]
        for key in keys_to_remove:
            del self.coeffs[key]

    def __call__(self, x):
        """Evaluates the polynomial at a given value x."""
        result = 0.0
        for power, coeff in self.coeffs.items():
            result += coeff * (x ** power)
        return result

    def __str__(self):
        """Returns a human-readable string representation of the polynomial."""
        if not self.coeffs:
            return "0"

        terms = []
        # Sort powers in descending order
        sorted_powers = sorted(self.coeffs.keys(), reverse=True)

        for power in sorted_powers:
            coeff = self.coeffs[power]
            if coeff == 0:
                continue

            coeff_abs_str = f"{abs(coeff):g}" # :g removes trailing .0 for integers
            
            if power == 0:
                term_str = coeff_abs_str
            elif power == 1:
                term_str = "x" if abs(coeff) == 1 else f"{coeff_abs_str}x"
            else:
                term_str = f"x^{power}" if abs(coeff) == 1 else f"{coeff_abs_str}x^{power}"
            
            if coeff < 0:
                terms.append(f"- {term_str}")
            else:
                terms.append(f"+ {term_str}")
        
        # Join terms, handling the first term's sign
        result_str = " ".join(terms)
        if result_str.startswith("+ "):
            result_str = result_str[2:] # Remove leading "+ "
        
        return result_str.strip()

    def __repr__(self):
        """Returns a developer-friendly string representation."""
        return f"Polynomial({dict(self.coeffs)})"

    def __add__(self, other):
        """Adds two polynomials."""
        if not isinstance(other, Polynomial):
            return NotImplemented
        
        new_coeffs = collections.defaultdict(float)
        for power, coeff in self.coeffs.items():
            new_coeffs[power] += coeff
        for power, coeff in other.coeffs.items():
            new_coeffs[power] += coeff
        
        return Polynomial(new_coeffs)

    def __sub__(self, other):
        """Subtracts one polynomial from another."""
        if not isinstance(other, Polynomial):
            return NotImplemented
        
        new_coeffs = collections.defaultdict(float)
        for power, coeff in self.coeffs.items():
            new_coeffs[power] += coeff
        for power, coeff in other.coeffs.items():
            new_coeffs[power] -= coeff
        
        return Polynomial(new_coeffs)

    def __mul__(self, other):
        """Multiplies two polynomials."""
        if not isinstance(other, Polynomial):
            return NotImplemented
        
        new_coeffs = collections.defaultdict(float)
        for p1, c1 in self.coeffs.items():
            for p2, c2 in other.coeffs.items():
                new_coeffs[p1 + p2] += c1 * c2
        
        return Polynomial(new_coeffs)

    def differentiate(self):
        """Returns a new Polynomial object representing the derivative."""
        new_coeffs = collections.defaultdict(float)
        for power, coeff in self.coeffs.items():
            if power > 0:
                new_coeffs[power - 1] = coeff * power
        
        return Polynomial(new_coeffs)

    def degree(self):
        """Returns the degree of the polynomial."""
        if not self.coeffs:
            return -1 # Convention for zero polynomial
        return max(self.coeffs.keys())

    def get_coefficient(self, power):
        """Returns the coefficient for a given power."""
        return self.coeffs.get(power, 0.0)