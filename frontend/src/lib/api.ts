import axios from "axios";
import {useEffect, useRef} from "react";

const API_URL = "https://api.fruitspace.one/v2/"

const api = axios.create({
    baseURL: API_URL,
})

export const gdps_top = async (page = 0) => {
    const {data} = await api.get(`fetch/gd/top?offset=${page}`)
    return data
}
export const gdps_get = async (srvid: string): Promise<GDPS | null> => {
    const {data} = await api.get<GDPS>(`fetch/gd/info/${srvid}`)
    return data.srvid ? data: null
}

export const gdps_repatch = async (srvid: string): Promise<Repatch> => {
    const {data} = await api.get<Repatch>(`repatch/gd/${srvid}`)
    return data
}

export type Repatch = {
    name: string,
    srvid: string,
    players: number,
    levels: number,
    icon: string,
    version: "2.2" | "2.1" | "2.0" | "1.9"
    recipe: string
}


export type GDPS = {
    srvid: string,
    plan: number,
    srv_name: string,
    owner_id: string,
    user_count: number,
    level_count: number,
    client_android_url: string,
    client_ios_url: string,
    client_windows_url: string,
    client_macoa_url: string,
    icon: string,
    description: string,
    text_align: number,
    discord: string,
    vk: string,
    version: "2.2" | "2.1" | "2.0" | "1.9",
    is_custom_textures: boolean,
    downloadpage_style: {
        accent: string,
        bg: string,
        variant: string
    }
}


export const useInterval = (callback: ()=>void, delay: number) => {

    const savedCallback = useRef<typeof callback>()

    useEffect(() => {
        savedCallback.current = callback
    }, [callback])


    useEffect(() => {
        function tick() {
            savedCallback.current?.()
        }
        if (delay !== null) {
            const id = setInterval(tick, delay)
            return () => clearInterval(id)
        }
    }, [delay])
}
