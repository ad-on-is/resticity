module.exports = {
	plugins: [require('@tailwindcss/typography'), require('daisyui')],
	theme: {
		extend: {
			colors: {
				'primary': '#cba6f7',
				'secondary': '#74c7ec',
				'accent': '#94e2d5',
				'neutral': '#313244',
				'base-100': '#1e1e2e',
				'info': '#74c7ec',
				'success': '#a6e3a1',
				'warning': '#f9e2af',
				'error': '#f38ba8',
			},
		},
	},
	daisyui: {
		themes: [
			'dark',
			{
				resticity: {
					'primary': '#cba6f7',
					'secondary': '#74c7ec',
					'accent': '#94e2d5',
					'neutral': '#313244',
					'base-100': '#1e1e2e',
					'info': '#74c7ec',
					'success': '#a6e3a1',
					'warning': '#f9e2af',
					'error': '#f38ba8',
				},
			},
		],
	},
}
