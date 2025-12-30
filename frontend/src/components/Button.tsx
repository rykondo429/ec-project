import React from 'react';
import Link from 'next/link';

export const Button: React.FC<
  React.ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: 'primary' | 'secondary' | 'outline' | 'danger';
    size?: 'sm' | 'md' | 'lg';
    asLink?: boolean;
    href?: string;
  }
> = ({ variant = 'primary', size = 'md', asLink, href, className, children, ...props }) => {
  const baseStyles = 'font-medium rounded transition-colors duration-200 flex items-center justify-center gap-2';

  const variantStyles = {
    primary: 'bg-primary text-white hover:bg-blue-600 disabled:bg-gray-300',
    secondary: 'bg-secondary text-white hover:bg-green-600 disabled:bg-gray-300',
    outline: 'border border-primary text-primary hover:bg-blue-50 disabled:border-gray-300',
    danger: 'bg-danger text-white hover:bg-red-600 disabled:bg-gray-300',
  };

  const sizeStyles = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg',
  };

  const combinedClassName = `${baseStyles} ${variantStyles[variant]} ${sizeStyles[size]} ${className || ''}`;

  if (asLink && href) {
    return (
      <Link href={href} className={combinedClassName}>
        {children}
      </Link>
    );
  }

  return (
    <button className={combinedClassName} {...props}>
      {children}
    </button>
  );
};
