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
import { Badge } from "@/components/ui/badge"
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
dayjs.extend(relativeTime)

interface TableProps {
    events: Event[]
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

export const EventsTable: FC<TableProps> = ({events}) => {
    return (
        <ScrollArea className="rounded-md border">
            <Table>
                <TableCaption>A list of your session events</TableCaption>
                <TableHeader>
                    <TableRow>
                        <TableHead className="w-[200px]">From</TableHead>
                        <TableHead className={"w-[400px]"}>Data</TableHead>
                        <TableHead className="text-center w-[80px]">Status</TableHead>
                        <TableHead className="text-right w-[400px]">Data</TableHead>
                        <TableHead className={"text-right w-[200px]"}>To</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {events.map((event) => (
                        <TableRow key={event.id}>
                            <TableCell className="font-medium">
                                <Badge className={"mb-2"} variant="secondary">{event.from.provider}</Badge>
                                 <br/>
                                {event.from.name}
                                <br/>
                                <Badge className={"mt-2"} variant="secondary">{dayjs(event.collected_at).fromNow()}</Badge>
                            </TableCell>
                            <TableCell>
                                <pre>{JSON.stringify(JSON.parse(event.from.data), null, 2)}</pre>
                            </TableCell>
                            <TableCell className="text-center border">
                                {statusIcon(event.state)} <br/>
                            </TableCell>
                            <TableCell >
                                <div className={"flex justify-end h-100"}>
                                    {
                                        event.state === "failed" ?
                                            <p className={"text-right"}>{event.message}</p>:
                                            <pre>{JSON.stringify(JSON.parse(event.to.data), null, 2)}</pre>
                                    }
                                </div>
                            </TableCell>
                            <TableCell className="font-medium text-right">
                                <Badge className={"mb-2"} variant="secondary">{event.to.provider}</Badge>
                                <br/>
                                {event.to.name}
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </ScrollArea>
    );
}
