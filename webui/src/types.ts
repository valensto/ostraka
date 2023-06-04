export interface Workflow {
    name: string;
    slug: string;
    nb_inputs: number;
    nb_outputs: number;
    events: Events;
}

export interface Event {
    workflow_slug: string;
    action: "received" | "sent";
    notifier: string;
    data: string;
    state: "succeed" | "failed";
    message: string;
}

export interface Events {
    received: Event[];
    sent: Event[];
}