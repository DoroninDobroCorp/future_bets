import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { Menu } from "@/components/Menu/Menu";


function Home({ }): JSX.Element {

	return (
		<>
			<Menu />
		</>
	)
}

export default withLayout(Home);