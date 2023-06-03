import {Workflow} from "../types";

const API_URL = process.env.NODE_ENV === 'production' ? '/' : 'http://localhost:4000';

export const getWorkflows = async (): Promise<Workflow[]> => {
    const res = await fetch(`${API_URL}/workflows`);
    return await res.json();
}

export function syncNotifications(onMessage: (event: MessageEvent) => void): EventSource | null {
    const eventSource = new EventSource(`${API_URL}/notifications`);
    eventSource.onmessage = onMessage;
    return eventSource;
}