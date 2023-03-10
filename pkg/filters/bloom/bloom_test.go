package bloomfilter

import (
	"testing"
	"unsafe"
)

func Test_hashValue(t *testing.T) {
	type args struct {
		value *[]byte
	}
	arr := []byte("hello")

	tests := []struct {
		name      string
		args      args
		wantHash1 uint
		wantHash2 uint
	}{
		{
			"test1",
			args{
				(*[]byte)(unsafe.Pointer(&arr)),
			},
			0x0,
			0x0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash1, gotHash2 := hashValue(tt.args.value)
			if gotHash1 != tt.wantHash1 {
				t.Errorf("hashValue() gotHash1 = %v, want %v", gotHash1, tt.wantHash1)
			}
			if gotHash2 != tt.wantHash2 {
				t.Errorf("hashValue() gotHash2 = %v, want %v", gotHash2, tt.wantHash2)
			}
		})
	}
}
