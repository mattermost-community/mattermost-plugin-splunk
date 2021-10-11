# Include custom targets and environment variables here

.PHONY: splunk
splunk:
	cd dev && docker-compose up
