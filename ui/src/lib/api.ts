import {Workflow} from "@/types";

const BASE_URL = process.env.NODE_ENV !== 'production' ? 'http://localhost:4000' : "";

export const getWorkflows = async (): Promise<Workflow[]> => {
    const res = await fetch(`${BASE_URL}/webui/workflows`);
    const workflows = await res.json();

    return workflows.map((workflow: Workflow) => ({
        ...workflow,
        events: []
    }));
}

export function syncEvents(onMessage: (event: MessageEvent) => void): EventSource | null {
    const eventSource = new EventSource(`${BASE_URL}/webui/consumes?token=2dc7929e5b589cb7861bcae19e13ad96`);
    eventSource.onmessage = onMessage;
    return eventSource;
}