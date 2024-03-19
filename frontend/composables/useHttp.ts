type FetchMethod = 'get' | 'post' | 'put' | 'delete' | 'patch' | 'head' | 'options'

export default class HttpClient {
	public static get = async (url: string, query: any = {}, notify: false | { title: string; text: string; type?: string } = false) =>
		await this.doFetch(url, { method: 'get', query }, notify)
	public static del = async (url: string, query: any = {}, notify: false | { title: string; text: string; type?: string } = false) =>
		await this.doFetch(url, { method: 'delete', query }, notify)
	public static post = async (url: string, data: any = {}, query: any = {}, notify: false | { title: string; text: string; type?: string } = false) =>
		await this.doFetch(url, { method: 'post', query, body: data }, notify)
	public static put = async (url: string, data: any, query: any = {}, notify: false | { title: string; text: string; type?: string } = false) =>
		await this.doFetch(url, { method: 'put', query, body: data }, notify)

	public static doFetch = async (url: string, opts: { method: FetchMethod; body?: any; query?: any }, notify: false | { title: string; text: string; type?: string } = false) => {
		const getUrl = (): string => {
			const url = useRequestURL()
			return url.protocol === 'wails:' || url.host.includes('wails.localhost') ? 'http://localhost:11278' : `${url.protocol}//${url.host}`
		}
		const baseUrl = `${getUrl()}/api`
		try {
			const res = await $fetch.raw(`${baseUrl}${url}`, {
				method: opts.method,
				body: JSON.stringify(opts.body),
				query: opts.query,
				headers: {
					'content-type': 'application/json',
				},
			})

			if (notify) {
				useToast().add({ id: url, title: notify.title, description: notify.text, icon: 'i-heroicons-eye' })
			}
			return res._data
		} catch (e: any) {
			console.error(e)
			useLogs().setServerError(e)
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
	}
}
