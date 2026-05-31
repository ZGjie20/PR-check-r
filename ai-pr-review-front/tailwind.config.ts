import type { Config } from 'tailwindcss';

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Noto Sans SC"', 'system-ui', 'sans-serif'],
      },
      colors: {
        sakura: {
          50: '#fff5f8',
          100: '#ffe4ec',
          200: '#ffc9db',
          300: '#ff9ec0',
          400: '#ff6ba3',
          500: '#f472b6',
        },
        dream: {
          50: '#f5f3ff',
          100: '#ede9fe',
          200: '#ddd6fe',
          300: '#c4b5fd',
          400: '#a78bfa',
          500: '#8b5cf6',
          600: '#7c3aed',
          700: '#6d28d9',
          800: '#5b21b6',
          900: '#4c1d95',
        },
      },
      animation: {
        float: 'float 6s ease-in-out infinite',
        'float-delayed': 'float 8s ease-in-out 2s infinite',
        'float-slow': 'float 10s ease-in-out 1s infinite',
        shimmer: 'shimmer 3s linear infinite',
        'pulse-glow': 'pulse-glow 2s ease-in-out infinite',
        'spin-slow': 'spin 8s linear infinite',
        twinkle: 'twinkle 3s ease-in-out infinite',
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0px) translateX(0px)' },
          '50%': { transform: 'translateY(-20px) translateX(10px)' },
        },
        shimmer: {
          '0%': { backgroundPosition: '200% center' },
          '100%': { backgroundPosition: '-200% center' },
        },
        'pulse-glow': {
          '0%, 100%': { opacity: '1', boxShadow: '0 0 8px rgba(167, 139, 250, 0.4)' },
          '50%': { opacity: '0.8', boxShadow: '0 0 20px rgba(244, 114, 182, 0.6)' },
        },
        twinkle: {
          '0%, 100%': { opacity: '0.3', transform: 'scale(1)' },
          '50%': { opacity: '1', transform: 'scale(1.2)' },
        },
      },
      backdropBlur: {
        xs: '2px',
      },
    },
  },
  plugins: [],
} satisfies Config;
