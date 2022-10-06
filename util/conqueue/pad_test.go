package conqueue

//
//type NoPad struct {
//	a uint64
//	b uint64
//	c uint64
//}
//
//func (np *NoPad) Increase() {
//	atomic.AddUint64(&np.a, 1)
//	atomic.AddUint64(&np.b, 1)
//	atomic.AddUint64(&np.c, 1)
//}
//
//type Pad struct {
//	a   uint64
//	_p1 [8]uint64
//	b   uint64
//	_p2 [8]uint64
//	c   uint64
//	_p3 [8]uint64
//}
//
//func (p *Pad) Increase() {
//	atomic.AddUint64(&p.a, 1)
//	atomic.AddUint64(&p.b, 1)
//	atomic.AddUint64(&p.c, 1)
//}
//func BenchmarkPad_Increase(b *testing.B) {
//	pad := &Pad{}
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			pad.Increase()
//		}
//	})
//}
//func BenchmarkNoPad_Increase(b *testing.B) {
//	nopad := &NoPad{}
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			nopad.Increase()
//		}
//	})
//}
