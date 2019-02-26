package main

const JobPath = "/api/jobs"

func (s *server) routes() {
	s.mux.HandleFunc(JobPath, s.getJobs).Methods("GET")
	s.mux.HandleFunc(JobPath+"/{id}", s.getJob).Methods("GET")
	s.mux.HandleFunc(JobPath, s.createJobs).Methods("POST")
	s.mux.HandleFunc(JobPath+"/{id}", s.updateJob).Methods("PATCH")
	s.mux.HandleFunc(JobPath+"/{id}", s.deleteJob).Methods("DELETE")
}
