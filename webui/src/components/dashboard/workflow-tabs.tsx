import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/components/ui/tabs";
import {EventsTable} from "./events-table";
import {Workflow} from "../../types";
import {FC} from "react";

interface WorkflowTabsProps {
    workflows: Workflow[]
}

export const WorkflowTabs: FC<WorkflowTabsProps> = ({workflows}) => {
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
                        <Card className="col-span-8">
                            <CardHeader>
                                <CardTitle>Inputs</CardTitle>
                            </CardHeader>
                            <CardContent className="pl-2">
                                <EventsTable events={workflow.events}/>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>
            ))}
        </Tabs>
    );
};
