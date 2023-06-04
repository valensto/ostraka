import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/components/ui/tabs";
import {InputTable, OutputTable} from "./events-table";
import {Workflow} from "../../types";
import {FC, useState} from "react";

interface WorkflowTabsProps {
    workflows: Workflow[]
}

export const WorkflowTabs: FC<WorkflowTabsProps> = ({workflows}) => {
    const [selectedRow, setSelectedRow] = useState<string>()

    if (workflows.length === 0) {
        return (
            <p className="text-sm text-muted-foreground">
                You don't have any workflow yet. Create one to get started.
            </p>
        )
    }

    return (
        <Tabs defaultValue={workflows[0].slug} className="space-y-4">
            <TabsList>
                {workflows.map((workflow) => (
                    <TabsTrigger key={workflow.slug} value={workflow.slug}>{workflow.name}</TabsTrigger>
                ))}
            </TabsList>
            {workflows.map((workflow) => (
                <TabsContent key={workflow.slug} value={workflow.slug} className="space-y-4">
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-8">
                        <Card className="col-span-4">
                            <CardHeader>
                                <CardTitle>Inputs</CardTitle>
                            </CardHeader>
                            <CardContent className="pl-2">
                                <InputTable selectRow={setSelectedRow} selectedRow={selectedRow} events={workflow.events.received}/>
                            </CardContent>
                        </Card>
                        <Card className="col-span-4">
                            <CardHeader>
                                <CardTitle>Outputs</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <OutputTable selectRow={setSelectedRow} selectedRow={selectedRow} events={workflow.events.sent}/>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>
            ))}
        </Tabs>
    );
};
