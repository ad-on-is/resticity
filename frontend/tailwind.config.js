module.exports = {
	plugins: [require('@tailwindcss/typography'), require('daisyui')],
	theme: {
		extend: {
			colors: {
				resticity: {
					50: '#cdd6f4',
					100: '#bac2de',
					200: '#a6adc8',
					300: '#9399b2',
					400: '#7f849c',
					500: '#6c7086',
					600: '#585b70',
					700: '#45475a',
					800: '#313244',
					900: '#1e1e2e',
					950: '#11111b',
				},
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
