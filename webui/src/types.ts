export interface Workflow {
    name: string;
    nb_inputs: number;
    nb_outputs: number;
}

export interface Notification {
    workflow: string;
    action: "received" | "sent";
    notifier: string;
    event: string;
    state: "succeed" | "failed";
    message: string;
}

export interface Notifications {
    Inputs: Notification[];
    Outputs: Notification[];
}