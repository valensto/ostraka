import {Workflow} from "@/types";

const BASE_URL = process.env.NODE_ENV !== 'production' ? 'http://localhost:4000' : "";

export const getWorkflows = async (): Promise<Workflow[]> => {
    const res = await fetch(`${BASE_URL}/webui/workflows`);
    const workflows = await res.json();

    return workflows.map((workflow: Workflow) => ({
        ...workflow,
        events: {
            received: [] as Event[],
            sent: [] as Event[],
        }
    }));
}

export function syncEvents(onMessage: (event: MessageEvent) => void): EventSource | null {
    const eventSource = new EventSource(`${BASE_URL}/webui/consume`);
    eventSource.onmessage = onMessage;
    return eventSource;
}