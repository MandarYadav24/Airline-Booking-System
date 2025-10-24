package main

import (
	"context"
	"log"
	"net"

)

func main() {
	l := logger.Init()
	l.Info().Msg("Starting booking service")

	ctx := context.Background()
	pool,err := db.NewPool(ctx, "postgres://airline:airlinepwd@localhost:5432/airlinedb")
	if err != nil {
		l.Fatal().Err(err).Msg("pg")
	}
	_ = pool
	writer := kafka.NewWriter("kafka:9092", "bookings")
	_ = writer

	lis, err := net.Listen("tcp", ":50052")
	if err != nil { 
		log.Fatalf("failed to listen: %v", err) 
	}
	s := grpc.NewServer()
	bookingpb.RegisterBookingServiceServer(s, &BookingServer{ /* deps */ })
	l.Info().Msg("booking service listening :50052")
	if err := s.Serve(lis); err != nil { 
		log.Fatalf("serve err: %v", err) 
	}
}
