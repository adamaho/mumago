<script lang="ts">
	import * as jsonpatch from 'fast-json-patch';
	import { onMount } from 'svelte';

	let form: HTMLFormElement;
	let data: any;
	let firstRequest = false;

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		const data = new FormData(form);
		const task = data.get('task');
		await fetch(`https://localhost:3000/todos/bar/${task}`, {
			method: 'POST',
			body: ''
		});
	}

	async function getTodos() {
		const todos = await fetch(`https://localhost:3000/todos/bar`, {
			headers: {
				'X-Muma-Stream': 'true'
			}
		});
		const reader = todos.body?.getReader();

		// eslint-disable-next-line no-constant-condition
		while (true) {
			if (!reader) {
				console.log('no reader?');
				break;
			}

			const { done, value } = await reader.read();
			if (done) {
				console.log('done?');
				break;
			}

			const decoder = new TextDecoder('utf-8');
			const str = decoder.decode(value.buffer);

			if (!firstRequest) {
				data = JSON.parse(str);
				firstRequest = true;
			} else {
				const patch = str.replace('\n', '');
				data = jsonpatch.applyPatch(data, JSON.parse(patch)).newDocument;
			}
		}
	}

	onMount(() => {
		getTodos();
	});
</script>

<form on:submit={handleSubmit} bind:this={form}>
	<label>
		Task
		<input name="task" />
	</label>
	<button>Add</button>
</form>
{#if data?.data}
	{#each data.data as d}
		<div>{d.task}</div>
	{/each}
{/if}
