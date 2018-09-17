const combat = window.combat = window.combat || {};
const TR_HAS_TRIES_CLASS = 'has-behaviour';
const TR_ACTIVE_CLASS = 'active';

combat.renderTable = function($target, logs) {
	const {createTag, showTries} = this;
	const $tries = createTag('div', {class: 'col-6'});
	const $tbody = createTag('tbody', {
		children: Object.keys(logs).map((key, index) => {
			const {state, title, tries} = logs[key];
			const hasBehaviour = tries && tries.length > 0;
			const $tr = createTag('tr', {
				class: hasBehaviour ? TR_HAS_TRIES_CLASS : '',
				children: [
					createTag('th', {scope: 'row', children: 1 + index}),
					createTag('td', {children:
						createTag('div', {class: `icon icon--${state}`})
					}),
					createTag('td', {children: title})
				]
			});

			if (hasBehaviour) {
				$tr.addEventListener('click', ({target}) => {
					const $prev = $tbody.querySelector(`.${TR_ACTIVE_CLASS}`);
					const parsedClassName = ` ${TR_ACTIVE_CLASS}`;

					if ($prev) {
						$prev.className = $prev.className.replace(parsedClassName, '')
					}

					target.closest('tr').className += parsedClassName;
					showTries($tries, tries);
				}, false);
			}

			return $tr;
		})
	});

	$target.append(
		createTag('div', {class: 'container', children:
			createTag('div', {class: 'row', children: [
				createTag('div', {class: 'col-6', children:
					createTag('table', {class: 'table table-hover', children: [
						createTag('thead', {class: 'thead-light', children:
							createTag('tr', {class: 'thead-light', children: [
								createTag('th', {cls: 'col', scope: 'col', children: '#'}),
								createTag('th', {css: 'col', scope: 'col', children: 'Status'}),
								createTag('th', {cas: 'col-10', scope: 'col', children: 'Test name'})
							]})
						}),
						$tbody
					]})
				}),
				$tries
			]})
		})
	);
}
