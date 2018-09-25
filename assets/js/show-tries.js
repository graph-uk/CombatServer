const combat = window.combat = window.combat || {};
const NAVIGATION_LIST_CLASS = 'log-navigation';
const NAVIGATION_TITLE_CLASS = 'log-navigation__title';
const NAVIGATION_ITEM_CLASS = 'log-navigation__item';
const NAVIGATION_ITEM_ACTIVE_CLASS = 'log-navigation__item--active';
const LOG_CONTAINER_CLASS = 'log';

const renderTriesNavigation = (tries, $tryPlaceholder) => {
	const {createTag, showTryDetails} = combat;
	const $buttons = tries.map((item, index) => {
		const $el = createTag('div', {
			class: `${NAVIGATION_ITEM_CLASS}${index === 0 ? ' ' + NAVIGATION_ITEM_ACTIVE_CLASS : ''}`,
			children: 1 + index
		});

		$el.addEventListener('click', ({target}) => showTryDetails(item, $tryPlaceholder, target), false);

		return $el;
	});

	return {
		'$triesNavigation': createTag('div', {class: NAVIGATION_LIST_CLASS, children: [
			createTag('div', {class: NAVIGATION_TITLE_CLASS, children: 'Logs:'}),
			...$buttons
		]}),
		$buttons
	};
}

combat.showTries = ($el, tries) => {
	const {createTag, showTryDetails} = combat;
	const $tryPlaceholder = createTag('div', {className: LOG_CONTAINER_CLASS});
	const {$triesNavigation, $buttons} = renderTriesNavigation(tries, $tryPlaceholder);

	$el.innerHTML = '';
	$el.append(
		$triesNavigation,
		$tryPlaceholder
	);

	showTryDetails(tries[0], $tryPlaceholder);
}
