import sys

def main():
    """
    Пример ввода и вывода числа n, где -10^9 < n < 10^9:
    n = int(input())
    print(n)
    """
    input_data = sys.stdin.read().split()
    n, T = map(int, [input_data[0], input_data[1]])

    timeline = [0] * (T + 1)
    print(timeline)
    count = []
    index = 2
    for line in range(n):
        a, f, s = map(int, [input_data[index], input_data[index+1], input_data[index+2]])
        a = max(a, 1)

        timeline[a] += s
        timeline[f] -= s
        index += 3

    del input_data
    now = 0
    maxi = 0
    for p in timeline:
        now += p
        maxi = max(maxi, now)

    print(maxi)
    return 0


if __name__ == '__main__':
    main()
