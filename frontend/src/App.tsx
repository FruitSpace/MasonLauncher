import {ReactNode, useEffect, useState} from 'react';

import BannerGD from "./assets/images/BannerGD.png"
import {GDPS, gdps_get, useInterval} from "./lib/api";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCirclePlay, faPlay, faPlusCircle, faRefresh, faUser} from "@fortawesome/free-solid-svg-icons";
import {ListServers, Patch, Read} from "../wailsjs/go/main/App";
import AnsiToHtml from "ansi-to-html"
import clsx from "clsx";
import {Dialog} from "radix-ui";


const VERSION = "1.0"


const AnsiConverter = new AnsiToHtml({
    colors: {
        34: "#0d6efd"
    }
})

const getServers = async (serverList: string[]) => {
    let mop: GDPS[] = []
    for (const srv of serverList) {
        const data = await gdps_get(srv)
        data && mop.push(data!)
    }
    return mop
}

function App() {
    const [clearance, setClearance] = useState<boolean>(false)
    const [srv, setSrv] = useState<GDPS>()
    const [serverList, setServerList] = useState<string[]>([])
    const [serversInfo, setServersinfo] = useState<GDPS[]>([])
    const [lockedState, setLockedState] = useState<boolean>(false)
    const [srvid, setSrvid] = useState<string>()
    const [dialogOpen, setDialogOpen] = useState(false)
    const [srvidToInstall,setSrvidToInstall] = useState<string>("")


    useEffect(() => {
        ListServers().then(v => {
            setServerList(v)
            srvid || setSrvid(v[0])
        })
    }, []);


    useEffect(
        () => {
            getServers(serverList).then(v => setServersinfo(v))
        },
        [serverList]
    )


    useEffect(() => {
        srvid && gdps_get(srvid).then(r => setSrv(r || undefined))
    }, [srvid]);



    return (
        <Dialog.Root open={dialogOpen}>
            <div className="z-[9999] fixed left-0 top-0 w-full h-8 draggable"></div>
            <div className="flex h-screen">
                <div className="h-full w-1/3 border-r-1 border-white border-opacity-25 p-4 pt-8 flex flex-col gap-4">
                    <p className="mx-auto text-gray-300 text-sm">Установленные серверы</p>
                    <div className="flex flex-col gap-4 overflow-y-scroll">
                        {serversInfo.map(srv => {
                            return <div
                                className={clsx("flex items-center gap-2 rounded-xl border-1 border-white border-opacity-25 cursor-pointer", srvid == srv.srvid && "bg-white bg-opacity-25")}
                                onClick={() => setSrvid(srv.srvid)}>
                                <img src={srv.icon} className="aspect-square w-16 rounded-xl"/>
                                <div className="">
                                    <p className="text ">{srv.srv_name}</p>
                                    <p className="text-xs text-slate-200">{srv.srvid}</p>
                                </div>
                            </div>
                        })}
                        <a href="#"
                           className="text-sm cursor-pointer hover:bg-white hover:bg-opacity-25 rounded-xl flex items-center justify-center gap-2 p-4"
                           onClick={() => setDialogOpen(true)}>
                            <FontAwesomeIcon icon={faPlusCircle}/> Добавить
                        </a>
                    </div>

                    <span className="mt-auto text-xs opacity-50">
                    Mason Launcher v{VERSION} by M41den
                </span>
                </div>
                <div className="h-full flex-1 bg-sidebar flex flex-col">
                    <div className="relative">
                        <img src={BannerGD}/>
                        <div
                            className="absolute top-0 left-0 w-full h-full bg-gradient-to-b from-transparent to-sidebar"></div>
                    </div>
                    <div className="flex flex-col relative z-50 mx-8 -mt-20 bg-black bg-opacity-75 rounded-xl">
                        <div
                            className=" flex gap-4 bg-sidebar border-white border-1 border-opacity-25 rounded-xl p-4">
                            <img src={srv?.icon} className="aspect-square w-24 rounded-lg"/>
                            <div className="">
                                <p className="text-2xl text-nowrap text-ellipsis">
                                    {srv?.srv_name || "Выберите сервер"}
                                </p>
                                <p className="text-slate-200 text-sm flex items-center gap-2">
                                    <span>
                                        <FontAwesomeIcon icon={faUser}/> {srv?.user_count || 0}
                                    </span>
                                    •
                                    <span>
                                        <FontAwesomeIcon icon={faCirclePlay}/> {srv?.level_count || 0}
                                    </span>
                                </p>
                            </div>
                        </div>
                        <div className="p-4 max-h-28 overflow-x-scroll overflow-y-scroll w-full">
                            <Log/>
                        </div>
                    </div>
                    <div className="flex gap-2 mx-auto mt-auto p-4 w-2/3">
                        <a href="#"
                           className={clsx(
                               "cursor-pointer px-4 py-2 bg-primary hover:bg-primary/80 rounded-lg flex-1 flex gap-2 items-center justify-center",
                               lockedState && "opacity-50 !cursor-default"
                           )}>
                            <FontAwesomeIcon icon={faPlay}/> Запустить
                        </a>
                        <a href="#"
                           className={clsx(
                               "cursor-pointer p-2 bg-success hover:bg-success/80 rounded-lg aspect-square h-10 flex gap-2 items-center justify-center",
                               lockedState && "opacity-50 !cursor-default"
                           )}
                           onClick={async () => {
                               if (!srv || lockedState) return
                               setLockedState(true)
                               await Patch(srv.srvid, srv.srv_name, srv.version).then(v => setLockedState(false))
                           }}>
                            <FontAwesomeIcon icon={faRefresh}/>
                        </a>
                    </div>
                </div>
            </div>
            <MasonDialog onClose={() => setDialogOpen(false)} okText="Установить" title="Добавить GDPS"
            onOk={async ()=>{
                const data = await gdps_get(srvidToInstall)
                if (!data) return
                setDialogOpen(false)
                setSrvid(srvidToInstall)
                setSrv(data)
                setLockedState(true)
                await Patch(data.srvid, data.srv_name, data.version).then(v => setLockedState(false))
            }}>
                <div className="flex justify-center">
                    <input placeholder="ID Сервера: XXXX"
                           value={srvidToInstall} autoComplete="off"
                           onChange={e=>setSrvidToInstall(e.target.value.replace(/[^a-zA-Z0-9]/, '').slice(0,4))}
                           className="px-3 py-1.5 rounded-lg outline-none text-sm bg-background"/>
                </div>
            </MasonDialog>
        </Dialog.Root>
    )
}

export default App


const Log = () => {
    const [log, setlog] = useState<string>("")

    useInterval(() => {
        Read().then(v => setlog((l) => `${l}${v}`))
    }, 200)
    return <p className="whitespace-pre font-mono text-xs text-slate-200" dangerouslySetInnerHTML={{
        __html: log ? AnsiConverter.toHtml(log.replaceAll("\n\n", "\n")) : "Здесь будет вывод консоли во время установки"
    }}></p>
}

const MasonDialog = (props: DialogProps) => {

    return <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 bg-black bg-opacity-40 data-[state=open]:animate-overlayShow z-50"/>
        <Dialog.Content
            className="z-[99] fixed left-1/2 top-1/2 max-h-[85vh] w-[90vw] max-w-[500px] -translate-x-1/2 -translate-y-1/2  p-4 shadow-[var(--shadow-6)] focus:outline-none data-[state=open]:animate-contentShow
            bg-sidebar rounded-xl border-1 border-white border-opacity-25">
            <Dialog.Title className="m-0 text-[17px] font-medium text-mauve12">
                {props.title}
            </Dialog.Title>
            {props.children}
            <div className="mt-[25px] flex gap-2 items-center justify-end">
                <a href="#" onClick={()=>props.onClose()} className="bg-red-500 px-3 py-1.5 rounded-md text-sm hover:bg-red-500/80">
                    {props.closeText || "Отмена"}
                </a>
                <a href="#" onClick={()=>props.onOk()} className="bg-primary px-3 py-1.5 rounded-md text-sm hover:bg-primary/80">
                    {props.okText || "ОК"}
                </a>
            </div>
        </Dialog.Content>
    </Dialog.Portal>

}

type DialogProps = {
    title: string,
    children: ReactNode,
    onClose: () => void,
    closeText?: string,
    onOk: () => void,
    okText: string
}