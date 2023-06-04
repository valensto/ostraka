import {useEffect, useState} from "react";
import {MainNav} from "@/components/dashboard/main-nav";
import {Metrics} from "@/components/dashboard/metrics";
import {WorkflowTabs} from "./components/dashboard/workflow-tabs";
import {Workflow, Event} from "./types";
import {getWorkflows, syncEvents} from "./lib/api";

export default function App() {
    const [workflows, setWorkflows] = useState<Workflow[]>([])

    useEffect(() => {
        getWorkflows().then((data) => {
            setWorkflows(data)
        }).catch((err) => {
            console.error(err)
        })
    }, [])

    useEffect(() => {
        const sync = syncEvents((event: MessageEvent) => {
            const data = JSON.parse(event.data) as Event

            setWorkflows(prevState => {
                return prevState.map((workflow) => {
                    if (workflow.slug === data.workflow_slug) {
                        const updatedEvents = {
                            ...workflow.events,
                            [data.action]: [
                                ...(workflow.events[data.action] || []),
                                data
                            ]
                        };
                        return {
                            ...workflow,
                            events: updatedEvents
                        };
                    }
                    return workflow;
                });
            });
        })

        return () => {
            if (sync) {
                sync.close()
            }
        }
    }, [])

    return (
        <div className="hidden flex-col md:flex">
            <div className="border-b">
                <MainNav className="mx-4"/>
            </div>
            <div className="flex-1 space-y-4 p-8 pt-6">
                <div className="flex items-center justify-between space-y-2">
                    <h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>
                </div>
                <Metrics workflows={workflows}/>
                <WorkflowTabs workflows={workflows}/>
            </div>
        </div>
    );
}
