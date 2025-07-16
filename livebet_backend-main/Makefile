MAKEFLAGS += j3

all: parse-pinnacle parse-lobbet parse-ladbrokes
	mkdir -p logs
	@echo "Run parser parrallel"

parse-pinnacle:
	cd parse_pinnacle && make >> ../logs/PinnacleLogs 2>&1

parse-lobbet:
	cd parse_lobbet && make >> ../logs/LobbetLogs 2>&1

parse-sansabet:
	cd parse_sansabet && make >> ../logs/SansabetLogs 2>&1

parse-betcenter:
	cd parse_betcenter && make >> ../logs/BetcenterLogs 2>&1

parse-ladbrokes:
	cd parse_ladbrokes && make >> ../logs/LadbrokesLogs 2>&1

parse-starcasino:
	cd parse_starcasino && make >> ../logs/StarCasinoLogs 2>&1

parse-unibet:
	cd parse_unibet && make >> ../logs/UnibetLogs 2>&1