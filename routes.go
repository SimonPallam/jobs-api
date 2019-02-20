package main

const JobPath = "/api/jobs"

func (s *server) routes() {
	s.mux.HandleFunc(JobPath, s.getJobs).Methods("GET")
	s.mux.HandleFunc(JobPath, s.createJobs).Methods("POST")
}
