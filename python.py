from line_profiler import LineProfiler
from collections import deque

class Solution:
    def maxVowels(self, s: str, k: int) -> int:
        last_k = deque(maxlen=k)
        [last_k.append(c in 'aeiou') for c in s[:k]]
        running_sum=sum(last_k)
        max=running_sum
        for c in s[k:]:
            value = c in 'aeiou'
            running_sum -= last_k.popleft() + value
            last_k.append(value)
            if running_sum>max:
                max=running_sum
        return max

# Create a profiler instance
profiler = LineProfiler()
# Wrap the function we want to profile
solution = Solution()
wrapped_func = profiler(solution.maxVowels)
# Run the function
wrapped_func("abciiidef" * 1000, 3)
# Print the stats
profiler.print_stats()