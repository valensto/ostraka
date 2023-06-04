import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {ScrollArea} from "@/components/ui/scroll-area";
import {FC} from "react";
import {Event} from "@/types";
import {ArrowRightFromLine, CircleSlash, Info} from "lucide-react";

interface TableProps {
    events: Event[]
    selectedRow?: string
    selectRow: (row: string) => void
}

const statusIcon = (status: string) => {
    switch (status) {
        case "succeed":
            return <span className="text-green-500"><ArrowRightFromLine style={{display: "inline-block"}}/></span>
        case "failed":
            return <span className="text-red-500"><CircleSlash style={{display: "inline-block"}}/></span>
        default:
            return <span className="text-gray-500"><Info style={{display: "inline-block"}}/></span>
    }
}

export const InputTable: FC<TableProps> = ({events, selectedRow, selectRow}) => {
    return (
        <ScrollArea className="h-[700px] rounded-md border">
            <Table>
                <TableCaption>A list of your session inputs events.</TableCaption>
                <TableHeader>
                    <TableRow>
                        <TableHead className="w-[150px]">From</TableHead>
                        <TableHead>Event</TableHead>
                        <TableHead>Message</TableHead>
                        <TableHead className="text-right">Status</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {events.map((event, k) => (
                        <TableRow
                            key={k}
                            onClick={() => selectRow(event.id)}
                            className={selectedRow === event.id ? "bg-clip-border border-4 border-violet-300 border-solid" : ""
                        }>
                            <TableCell className="font-medium">{event.notifier}</TableCell>
                            <TableCell>
                                <pre>{JSON.stringify(JSON.parse(event.data), null, 2)}</pre>
                            </TableCell>
                            <TableCell>{event.message}</TableCell>
                            <TableCell className="text-right">
                                {statusIcon(event.state)}
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </ScrollArea>
    );
}

export const OutputTable: FC<TableProps> = ({events, selectedRow, selectRow}) => {
    return (
        <ScrollArea className="h-[700px] rounded-md border">
            <Table>
                <TableCaption>A list of your session outputs events.</TableCaption>
                <TableHeader>
                    <TableRow>
                        <TableHead>To</TableHead>
                        <TableHead>Payload</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {events.map((event, k) => (
                        <TableRow
                            key={k}
                            onClick={() => selectRow(event.id)}
                            className={selectedRow === event.id ? "bg-clip-border border-4 border-violet-300 border-solid" : ""
                        }>
                            <TableCell className="font-medium">{event.notifier}</TableCell>
                            <TableCell>
                                <pre>
                                      {JSON.stringify(JSON.parse(event.data),null,2)}
                                </pre>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </ScrollArea>
    );
}