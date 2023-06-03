import {Workflow} from "../types";

const API_URL = process.env.NODE_ENV === 'production' ? '/' : 'http://localhost:4000';

export const getWorkflows = async (): Promise<Workflow[]> => {
    const res = await fetch(`${API_URL}/workflows`);
    const data = await res.json();
    return data.workflows;
}
