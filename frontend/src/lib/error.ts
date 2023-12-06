import { writable, type Writable } from 'svelte/store';

const error: Writable<string> = writable('');

export { error };
