package main

import (
	"bufio"
	"fmt"
	"os"
)

func readInt(reader *bufio.Reader) int {
	var res int
	b, err := reader.ReadByte()
	negative := 1
	for err == nil && (b == '\n' || b == ' ' || b == '\r') {
		b, err = reader.ReadByte()
	}

	if err != nil {
		return 0
	}

	if b == '-' {
		negative = -1
		b, err = reader.ReadByte()
	}

	for err == nil && b >= '0' && b <= '9' {
		res = res*10 + int(b-'0')
		b, err = reader.ReadByte()
	}
	res = res * negative
	return res
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)
	N := readInt(reader)

	left := make([]int, 0, N)
	right := make([]int, 0, N)
	up := make([]int, 0, N)
	down := make([]int, 0, N)

	for m := 0; m < 4; m++ {
		for i := 0; i < N; i++ {
			newInt := readInt(reader)
			if m == 0 {
				left = append(left, newInt)
			} else if m == 1 {
				right = append(right, newInt)
			} else if m == 2 {
				up = append(up, newInt)
			} else {
				down = append(down, newInt)
			}
		}
	}

	for globalIdx, left_idx := range left {
		right_val := right[globalIdx]
		if (left_idx == -1) != (right_val == -1) {
			fmt.Println("NO")
			return
		}
		if left_idx == -1 {
			continue
		}

		right_idx := N - right_val - 1
		if left_idx > right_idx {
			fmt.Println("NO")
			return
		}

		if left_idx != -1 && right_val != -1 {
			up_first_idx := up[left_idx]
			down_first_idx := N - down[left_idx] - 1

			up_second_idx := up[right_idx]
			down_second_idx := N - down[right_idx] - 1

			if up_first_idx == -1 || up_second_idx == -1 {
				fmt.Println("NO")
				return
			}

			if up_first_idx > down_first_idx || up_second_idx > down_second_idx {
				fmt.Println("NO")
				return
			}

			if up_first_idx > globalIdx || down_first_idx < globalIdx || up_second_idx > globalIdx || down_second_idx < globalIdx {
				fmt.Println("NO")
				return
			}

		}

	}

	for globalIdx, up_idx := range up {
		down_val := down[globalIdx]
		if (up_idx == -1) != (down_val == -1) {
			fmt.Println("NO")
			return
		}
		if up_idx == -1 {
			continue
		}

		down_idx := N - down_val - 1
		if up_idx > down_idx {
			fmt.Println("NO")
			return
		}
		if up_idx != -1 && down_val != -1 {
			left_first_idx := left[up_idx]
			right_first_idx := N - right[up_idx] - 1

			left_second_idx := left[down_idx]
			right_second_idx := N - right[down_idx] - 1

			if left_first_idx == -1 || left_second_idx == -1 {
				fmt.Println("NO")
				return
			}

			if left_first_idx > right_first_idx || left_second_idx > right_second_idx {
				fmt.Println("NO")
				return
			}

			if left_first_idx > globalIdx || right_first_idx < globalIdx || left_second_idx > globalIdx || right_second_idx < globalIdx {
				fmt.Println("NO")
				return
			}

		}
	}
	fmt.Println("YES")

}
