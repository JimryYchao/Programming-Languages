package gostd

/*
! ConstantTimeByteEq：x == y 返回 1，否则返回 0。
! ConstantTimeCompare：x == y 返回 1，否则返回 0。x 和 y 长度不一致，则立即返回 0。
! ConstantTimeCopy：v == 1，函数将 y 的内容复制到 x（等长切片）中。v == 0，则 x 保持不变。v 其他值行为未定义
! ConstantTimeEq：x == y 返回 1，否则返回 0。
! ConstantTimeLessOrEq：x <= y 返回 1，否则返回 0。如果 x 或 y 为负数或大于 2^31-1，则行为未定义
! ConstantTimeSelect：v == 1 返回 x; v == 0 返回 y; v 其他值行为未定义
! XORBytes：将所有 i < n = min(len(x),len(y)) 设置为 dst[i] = x[i] ^ y[i]，返回写入 dst 的字节数 n。
	如果 dst 的长度不大于 n，XOR 将 panic，而不向 dst 写入任何内容。
*/
