type FetchMethod = 'get' | 'post' | 'put' | 'delete' | 'patch' | 'head' | 'options'

export default class HttpClient {
	public static get = async (url: string, notify: false | { title: string; text: string; type?: string } = false) => await this.doFetch(url, { method: 'get' }, notify)
	public static del = async (url: string, notify: false | { title: string; text: string; type?: string } = false) => await this.doFetch(url, { method: 'delete' }, notify)
	public static post = async (url: string, data: any, notify: false | { title: string; text: string; type?: string } = false) =>
		await this.doFetch(url, { method: 'post', body: data }, notify)
	public static put = async (url: string, data: any, notify: false | { title: string; text: string; type?: string } = false) =>
		await this.doFetch(url, { method: 'put', body: data }, notify)

	public static doFetch = async (url: string, opts: { method: FetchMethod; body?: any }, notify: false | { title: string; text: string; type?: string } = false) => {
		// let baseUrl = useRuntimeConfig().public.apiURL
		// if (url.includes('clients') || url.includes('invoicesoroffers') || url.includes('timetracks') || url.includes('profile') || url.includes('auth')) {
		const baseUrl = 'http://localhost:11278/api'
		// }
		try {
			const res = await $fetch.raw(`${baseUrl}${url}`, {
				method: opts.method,
				body: JSON.stringify(opts.body),
				headers: {
					'content-type': 'application/json',
				},
			})

			if (notify) {
				useToast().add({ id: url, title: notify.title, description: notify.text, icon: 'i-heroicons-eye' })
			}
			return res._data
		} catch (e: any) {
			console.log(e)
			this.notifyError(e, notify)

			return e.data
		}
	}

	private static notifyError(e: any, notify: false | { title: string; text: string; type?: string } = false) {
		let title = 'Error'
		let message = 'Unexpected error occured'
		if (notify) {
			title = notify.title
		}
		if (e.data) {
			message = e.data
		}
		useToast().add({ title: title, description: message, icon: 'i-heroicons-exclamation-triangle', color: 'red' })

		// useNotification().notify({
		// 	title: title,
		// 	text: message,
		// 	type: 'error',
		// })
	}
}
