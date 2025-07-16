import { JSX, useEffect, useMemo, useState } from "react";
import styles from "./ClientPageComponent.module.css";
import { useUnit } from "effector-react";
import $bookmaker from "../../stores/BookmakerStore";
import { $pair, $pairCount, setPair, setPairCount } from "@/stores/PairStore";
import $filterBet from "../../stores/FilterBetStore";
import { $pairPrematch, $pairPrematchCount, setPairPrematchCount, setPairPrematch } from "@/stores/PairPrematchStore";
import { Match } from "@/components/Match/Match";
import { IPair, IPairOneOutcome } from "@/interfaces/Pair";
import { Calculator } from "@/components/Calculator/Calculator";
import { Status } from "@/components/Status/Status";
import { StatusError, StatusOK } from "@/helpers/TypeStatusByTime";
import { useWebSocket } from "@/websocket/websocket";
import { onError } from "@/websocket/error";
import { ClientMatchFilter } from "@/components/ClientMatchFilter/ClientMatchFilter";
import $sport from "../../stores/SportStore";
import {$workMode, EmployeesStore} from "@/stores/EmployeesStore";
import { IGeneralFilter } from "@/interfaces/GeneralFilter";

export const ClientPageComponent = ({ userId }: { userId: number }): JSX.Element => {
    const bookmakers = useUnit($bookmaker)
    const sports = useUnit($sport)
    const pairsLive = useUnit($pair)
    const pairsCount = useUnit($pairCount)
    const pairsPrematch = useUnit($pairPrematch)
    const pairsPrematchCount = useUnit($pairPrematchCount)
    const filterBet = useUnit($filterBet)
    const workMode = useUnit($workMode)

    const pairs: Map<string, IPair> = useMemo(() => {
        return new Map([...pairsLive, ...pairsPrematch]);
    }, [pairsLive, pairsPrematch]);

    const [selectGeneralFilter, setSelectGeneralFilter] = useState<IGeneralFilter>();
    const [selectKeyMatch, setSelectKeyMatch] = useState<string | null>(null);
    const [selectOutcome, setSelectOutcome] = useState<string | null>(null);
    const [selectMatch, setSelectMatch] = useState<IPairOneOutcome | null>(null);
    const [closeCalc, setCloseCalc] = useState<boolean>(true);
    const [sentPair, setSentPair] = useState<IPairOneOutcome | null>(null);

    /* eslint-disable @typescript-eslint/no-explicit-any */
    const parseAnalyzerMsg = (message: any) => {
        const msgs: IPair[] = JSON.parse(message.data)
        setPairCount(msgs ? msgs.length : 0)
        setPair(msgs)
    }

    const parseAnalyzerPreMsg = (message: any) => {
        const msgs: IPair[] = JSON.parse(message.data)
        setPairPrematchCount(msgs ? msgs.length : 0)
        setPairPrematch(msgs)
    }

    const [wsStatus, sendMessage] = useWebSocket(process.env.NEXT_PUBLIC_ANALYSER_WS!, //important
        parseAnalyzerMsg, onError);
    useEffect(() => {
        console.log("Webosocket analyzer status - " + wsStatus)
    }, [wsStatus]);
    const [wsStatusPre, sendMessagePre] = useWebSocket(process.env.NEXT_PUBLIC_ANALYSER_PRE_WS!,
        parseAnalyzerPreMsg, onError);
    useEffect(() => {
        console.log("Webosocket analyzer prematch status - " + wsStatusPre)
    }, [wsStatusPre])

    const handleMatch = (match: IPair, key: string) => {
        setSelectKeyMatch(key)
        setSelectOutcome(match.outcome[0].outcome)
    };

    useEffect(() => {
        const sendFilters = async () => {
            if (selectGeneralFilter && wsStatus === "open") {
                sendMessage(selectGeneralFilter);
                sendMessagePre(selectGeneralFilter);

                await EmployeesStore.sendFilters(userId, selectGeneralFilter);
            }
        };

        sendFilters();
    }, [selectGeneralFilter, workMode, userId, wsStatus, sendMessage, sendMessagePre]);

    useEffect(() => {
        if (pairs && selectKeyMatch && selectOutcome) {
            const select: IPair | undefined = pairs.get(selectKeyMatch)
            if (select) {
                const pair: IPair = JSON.parse(JSON.stringify(select))
                pair.outcome = pair.outcome.filter(out => {
                    return out.outcome == selectOutcome
                })
                if (pair.outcome.length > 0) {
                    const pairOneOutcome: IPairOneOutcome = {
                        first: pair.first,
                        second: pair.second,
                        outcome: pair.outcome[0],
                        isLive: pair.isLive,
                        sportName: pair.sportName,
                        createdAt: pair.createdAt
                    }
                    setSelectMatch(pairOneOutcome)
                }
            } else {
                // setSelectMatch(null)
            }
        } else {
            // setSelectMatch(null)
        }
    }, [pairs, selectKeyMatch, selectOutcome])

    useEffect(() => {
        if (!closeCalc) {
            setSelectKeyMatch(null)
            setSelectOutcome(null)
            setSelectMatch(null)
            setCloseCalc(true)
            setSentPair(null)
        }
    }, [closeCalc])

    return (
        <>
            <div className={styles.wrapper}>

                <div className={styles.webStatus}>
                    <Status title="Analyzer" status={(wsStatus == "open") ? StatusOK : StatusError} /> <p>{pairsCount}</p>
                    <Status title="AnalyzerPrematch" status={(wsStatusPre == "open") ? StatusOK : StatusError} /> <p>{pairsPrematchCount}</p>
                </div>

                {bookmakers && sports && <ClientMatchFilter bookmakers={bookmakers} sports={sports} setSelect={setSelectGeneralFilter} />}

                {selectMatch && <Calculator setSentPair={setSentPair} setClose={setCloseCalc} pair={sentPair || selectMatch} userId={userId} />}

                {workMode && pairs && [...pairs].sort((x, y) => y[1].outcome[0].roi - x[1].outcome[0].roi).filter(value => {
                    return !filterBet.includes(value[1].first.matchId) 
                }).map(value => {
                    return <Match key={value[0]} keyMatch={value[1].first.bookmaker + value[1].first.matchId + value[1].second.bookmaker + value[1].second.matchId + value[1].sportName} pair={value[1]} setSelectMatch={handleMatch} />
                })}
            </div>
        </>
    )
}