package domain

import "context"

// TxManager adalah abstraksi unit-of-work (persistence-convention.md §2) — dimiliki &
// dipanggil di Application Layer, BUKAN oleh repository. Implementasi konkret (GORM)
// ada di adapter layer, menyisipkan handle transaksi ke context supaya repository yang
// dipanggil di dalam fn membaca transaksi yang sama, bukan koneksi db dasar.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
