export interface Workflow {
    name: string;
    slug: string;
    nb_inputs: number;
    nb_outputs: number;
    events: Event[];
}

export interface Event {
    id: string;
    workflow_slug: string;
    from: Source;
    to: Source;
    state: "succeed" | "failed";
    message: string;
    collected_at: string;
}

interface Source {
    provider: string;
    name: string;
    data: string;
}