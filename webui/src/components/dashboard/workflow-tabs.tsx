import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { InputTable, OutputTable } from "./events-table";
import {Events, Workflow} from "../../types";
import {FC} from "react";

interface WorkflowTabsProps {
  workflows: Workflow[]
}

export const WorkflowTabs: FC<WorkflowTabsProps> = ({workflows}) => {
  return (
    <Tabs defaultValue="orders" className="space-y-4">
      <TabsList>
        {workflows.map((workflow) => (
            <TabsTrigger key={workflow.slug} value={workflow.slug}>{workflow.name}</TabsTrigger>
        ))}
      </TabsList>
      {workflows.map((workflow) => (
          <TabsContent key={workflow.slug} value={workflow.slug} className="space-y-4">
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-8">
                  <Card className="col-span-5">
                      <CardHeader>
                          <CardTitle>Inputs</CardTitle>
                      </CardHeader>
                      <CardContent className="pl-2">
                          <InputTable events={workflow.events.received} />
                      </CardContent>
                  </Card>
                  <Card className="col-span-3">
                      <CardHeader>
                          <CardTitle>Outputs</CardTitle>
                      </CardHeader>
                      <CardContent>
                          <OutputTable events={workflow.events.sent} />
                      </CardContent>
                  </Card>
              </div>
          </TabsContent>
      ))}
    </Tabs>
  );
};
