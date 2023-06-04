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
import {Event} from "../../types";

interface TableProps {
    events: Event[]
}

export const InputTable: FC<TableProps> = ({events}) => {
    return (
        <ScrollArea className="h-[700px] rounded-md border">
            <Table>
                <TableCaption>A list of your session inputs events.</TableCaption>
                <TableHeader>
                    <TableRow>
                        <TableHead>From</TableHead>
                        <TableHead>Payload</TableHead>
                        <TableHead>Status</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {events.map((event, k) => (
                        <TableRow key={k}>
                            <TableCell className="font-medium">{event.notifier}</TableCell>
                            <TableCell>{event.data}</TableCell>
                            <TableCell>
                                {event.state}
                                {event.message ? <br/> : ""}
                                {event.message}
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </ScrollArea>
    );
}

export const OutputTable: FC<TableProps> = ({events}) => {
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
                        <TableRow key={k}>
                            <TableCell className="font-medium">{event.notifier}</TableCell>
                            <TableCell>{event.data}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </ScrollArea>
    );
}