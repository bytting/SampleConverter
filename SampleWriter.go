package main

type SampleWriter interface {

	Write(s *Sample) error
	Close() error
}
