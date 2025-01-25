import {useCallback, useEffect, useMemo, useState} from 'react';

import BannerGD from "./assets/images/BannerGD.png"
import {GDPS, gdps_get, useInterval} from "./lib/api";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCirclePlay, faPlay, faRefresh, faUser} from "@fortawesome/free-solid-svg-icons";
import {Greet, ListServers, Patch, Read} from "../wailsjs/go/main/App";
import clsx from "clsx";

const VERSION = "1.0"

const getServers = async (serverList: string[])=>{
    let mop: GDPS[] = []
    for(const srv of serverList) {
        const data = await gdps_get(srv)
        data&&mop.push(data!)
    }
    return mop
}

function App() {
    const [srv, setSrv] = useState<GDPS>()
    const [serverList, setServerList] = useState<string[]>([])
    const [serversInfo, setServersinfo] = useState<GDPS[]>([])

    const [srvid, setSrvid] = useState<string>()

    console.log(srv,srvid,serverList,serversInfo)


    useEffect(() => {
        ListServers().then(v=>{
            setServerList(v)
            srvid||setSrvid(v[0])
        })
    }, []);


    useEffect(
        ()=>{
            getServers(serverList).then(v=>setServersinfo(v))
        },
        [serverList]
    )


    useEffect(() => {
        srvid&&gdps_get(srvid).then(r=>setSrv(r||undefined))
    }, [srvid]);


    return (
        <div className="flex h-screen">
            <div className="z-10 fixed left-0 top-0 w-full h-8 draggable"></div>
            <div className="h-full w-1/3 border-r-1 border-white border-opacity-25 p-4 pt-8 flex flex-col gap-4">
                <p className="mx-auto text-gray-300 text-sm">Установленные серверы</p>
                <div className="flex flex-col gap-4 overflow-y-scroll">
                    {serversInfo.map(srv=>{
                        return <div className={clsx("flex items-center gap-2 rounded-xl border-1 border-white border-opacity-25 cursor-pointer", srvid==srv.srvid && "bg-white bg-opacity-25")}
                                    onClick={()=>setSrvid(srv.srvid)}>
                            <img src={srv.icon} className="aspect-square w-16 rounded-lg" />
                            <div className="">
                                <p className="text ">{srv.srv_name}</p>
                                <p className="text-xs text-slate-200">{srv.srvid}</p>
                            </div>
                        </div>
                    })}
                </div>

                <span className="mt-auto text-xs opacity-50">
                    Mason Launcher v{VERSION} by M41den
                </span>
            </div>
            <div className="h-full flex-1 bg-sidebar flex flex-col">
                <div className="relative">
                    <img src={BannerGD}/>
                    <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-b from-transparent to-sidebar"></div>
                </div>
                <div className="flex flex-col relative z-50 mx-8 -mt-20 bg-black bg-opacity-75 rounded-xl">
                    <div
                        className=" flex gap-4 bg-sidebar border-white border-1 border-opacity-25 rounded-xl p-4">
                        <img src={srv?.icon} className="aspect-square w-24 rounded-lg"/>
                        <div className="">
                            <p className="text-2xl text-nowrap text-ellipsis">{srv?.srv_name}</p>
                            <p className="text-slate-200 text-sm">
                                <FontAwesomeIcon icon={faUser}/> {srv?.user_count} • <FontAwesomeIcon
                                icon={faCirclePlay}/> {srv?.level_count}
                            </p>
                        </div>
                    </div>
                    <div className="whitespace-pre font-mono p-4 text-xs text-slate-200 max-h-28 overflow-x-scroll overflow-y-scroll">
                        <Log />
                    </div>
                </div>
                <div className="flex gap-2 mx-auto mt-auto p-4 w-2/3">
                    <a href="#"
                       className="cursor-pointer px-4 py-2 bg-primary hover:bg-primary/80 rounded-lg flex-1 flex gap-2 items-center justify-center">
                        <FontAwesomeIcon icon={faPlay}/> Запустить
                    </a>
                    <a href="#"
                       className="cursor-pointer p-2 bg-success hover:bg-success/80 rounded-lg aspect-square h-10 flex gap-2 items-center justify-center"
                       onClick={() => {
                           srv && Patch(srv.srvid, srv.version)
                       }}>
                    <FontAwesomeIcon icon={faRefresh} />
                    </a>
                </div>
            </div>
        </div>
    )
}

export default App


const Log = () => {
    const [log, setlog] = useState<string>("")

    useInterval(()=>{
        console.log("poll")
        Read().then(v=>setlog((l)=>`${l}${v}`))
    }, 200)
    return <>
        {
            log||"Здесь будет вывод консоли во время установки"
        }
    </>
}