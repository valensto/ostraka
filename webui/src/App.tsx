import {useCallback, useEffect, useState} from "react";
import {MainNav} from "@/components/dashboard/main-nav";
import {Metrics} from "@/components/dashboard/metrics";
import {WorkflowTabs} from "./components/dashboard/workflow-tabs";
import {Notifications, Workflow} from "./types";
import {getWorkflows} from "./lib/api";

export default function App() {
    const [notifications, setNotifications] = useState<Notifications>({
        Inputs: [],
        Outputs: [],
    })
    const [workflows, setWorkflows] = useState<Workflow[]>([])

    useEffect(() => {
        getWorkflows().then((data) => {
            setWorkflows(data)
        }).catch((err) => {
            console.error(err)
        })
    }, [])

    return (
        <>
            <div className="hidden flex-col md:flex">
                <div className="border-b">
                    <MainNav className="mx-4"/>
                </div>
                <div className="flex-1 space-y-4 p-8 pt-6">
                    <div className="flex items-center justify-between space-y-2">
                        <h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>
                    </div>
                    <Metrics workflows={workflows} notifications={notifications}/>
                    <WorkflowTabs/>
                </div>
            </div>
        </>
    );
}
