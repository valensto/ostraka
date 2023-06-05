import {FC, useEffect, useState} from "react";
import {Activity, FileInput, FileOutput} from "lucide-react";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import {Workflow} from "@/types"

interface MetricsProps {
    workflows: Workflow[]
}

interface counter {
    Inputs: number
    Outputs: number
    InputsEvents: number
    OutputsEvents: number
}

export const Metrics: FC<MetricsProps> = ({workflows}) => {
    const [counter, setCounter] = useState<counter>({
        Inputs: 0,
        Outputs: 0,
        InputsEvents: 0,
        OutputsEvents: 0
    })

    useEffect(() => {
        if (workflows.length === 0) {
            return
        }

        let inputs = 0
        let outputs = 0

        workflows.forEach((workflow) => {
            inputs += workflow.nb_inputs
            outputs += workflow.nb_outputs
        })

        setCounter(prevState => {
            return {
                ...prevState,
                Inputs: inputs,
                Outputs: outputs
            }
        })
    }, [workflows])

    return (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">
                        Inputs registered
                    </CardTitle>
                    <FileInput className="h-4 w-4 text-muted-foreground"/>
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">{counter.Inputs}</div>
                </CardContent>
            </Card>
            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Input Events</CardTitle>
                    <Activity className="h-4 w-4 text-muted-foreground"/>
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">+{0}</div>
                    <p className="text-xs text-muted-foreground">for this session</p>
                </CardContent>
            </Card>
            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">
                        Outputs registered
                    </CardTitle>
                    <FileOutput className="h-4 w-4 text-muted-foreground"/>
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">{counter.Outputs}</div>
                </CardContent>
            </Card>
            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Output Events</CardTitle>
                    <Activity className="h-4 w-4 text-muted-foreground"/>
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">+{0}</div>
                    <p className="text-xs text-muted-foreground">for this session</p>
                </CardContent>
            </Card>
        </div>
    );
};
