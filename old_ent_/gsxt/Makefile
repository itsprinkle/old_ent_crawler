all:
	go build -o build/apps/master apps/grab/master.go
	go build -o build/apps/dispatcher apps/grab/dispatcher.go
	go build -o build/apps/worker apps/grab/worker.go
	go build -o build/apps/find apps/find/find.go
	go build -o build/apps/v1 apps/v1/v1.go

grab:
	go build -o build/apps/dispatcher apps/grab/dispatcher.go
	go build -o build/apps/worker apps/grab/worker.go

v1:
	go build -o build/apps/v1 apps/v1/v1.go

find:
	go build -o build/apps/find apps/find/find.go

clean:
	$(RM) build/apps/find
	$(RM) build/apps/worker
	$(RM) build/apps/dispatcher
	$(RM) build/apps/master
	$(RM) build/apps/v1

.PHONY: all clean find grab v1
