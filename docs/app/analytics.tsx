"use client";

import { usePathname } from "next/navigation";
import { useEffect } from "react";

import type { ReactNode } from "react";

export default function Analytics({ children }: { children: ReactNode }) {
	const pathname = usePathname();

	useEffect(() => {
		const handleRouteChange = (url) => {
			gtag("event", "page_view", {
				page_path: url,
			});
		};

		// Track the initial load
		handleRouteChange(pathname);

		// Listen to changes in the pathname
		return () => handleRouteChange(pathname);
	}, [pathname]);

	return children;
}
