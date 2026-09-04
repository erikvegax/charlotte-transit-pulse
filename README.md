# charlotte-transit-pulse
Charlotte Transit Pulse is an event-driven pipeline: a Go service polls CATS' GTFS-Realtime feed, publishes normalized events to Kafka, a consumer service joins those events against the static schedule and writes computed reliability data to MongoDB, and a Go API serves that data to a React dashboard.
