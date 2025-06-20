import math

def is_prime(n):
    if n < 2:
        return False
    if n == 2:
        return True
    if n % 2 == 0:
        return False
    i = 3
    while i * i <= n:
        if n % i == 0:
            return False
        i += 2
    return True

def prime_generator():
    num = 2
    while True:
        if is_prime(num):
            yield num
        num += 1

if __name__ == "__main__":
    prime_gen = prime_generator()
    # Generate the first 10 prime numbers
    for _ in range(10):
        print(next(prime_gen))

    # Generate prime numbers up to a certain limit (example)
    # print("\nPrimes up to 30:")
    # prime_gen_limited = prime_generator()
    # while True:
    #     p = next(prime_gen_limited)
    #     if p > 30:
    #         break
    #     print(p)

# Additional implementation at 2025-06-19 22:40:13
import math

class PrimeTools:
    def is_prime(self, n: int) -> bool:
        if n < 2:
            return False
        if n == 2:
            return True
        if n % 2 == 0:
            return False
        for i in range(3, int(math.sqrt(n)) + 1, 2):
            if n % i == 0:
                return False
        return True

    def generate_primes_up_to(self, limit: int) -> list[int]:
        if limit < 2:
            return []
        sieve = [True] * (limit + 1)
        sieve[0] = sieve[1] = False
        for p in range(2, int(math.sqrt(limit)) + 1):
            if sieve[p]:
                for multiple in range(p*p, limit + 1, p):
                    sieve[multiple] = False
        primes = [i for i, is_prime in enumerate(sieve) if is_prime]
        return primes

    def _prime_number_iterator(self):
        yield 2
        primes = [2]
        num = 3
        while True:
            is_p = True
            for p in primes:
                if p * p > num:
                    break
                if num % p == 0:
                    is_p = False
                    break
            if is_p:
                primes.append(num)
                yield num
            num += 2

    def generate_first_n_primes(self, count: int) -> list[int]:
        if count <= 0:
            return []
        primes = []
        prime_iter = self._prime_number_iterator()
        for _ in range(count):
            primes.append(next(prime_iter))
        return primes

    def get_primes_in_range(self, start: int, end: int) -> list[int]:
        if start > end:
            return []
        all_primes_up_to_end = self.generate_primes_up_to(end)
        primes_in_range = [p for p in all_primes_up_to_end if p >= start]
        return primes_in_range