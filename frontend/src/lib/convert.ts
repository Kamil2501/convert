import { error } from './error';

export async function convert(form: HTMLFormElement, convertOrder: string[]) {
	const formData = new FormData(form);
	const url = new URL(
		`http://localhost:8080/convert/${convertOrder[0]}to${convertOrder[1]}`
	);

	try {
		const response: Response = await fetch(url, {
			method: form.method,
			body: formData,
		});

		const data: Blob = await response.blob();

		return data;
	} catch (err) {
		error.set('Server error');
		console.log(err);
	}
}
