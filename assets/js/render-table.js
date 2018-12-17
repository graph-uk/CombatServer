const combat = window.combat = window.combat || {};
const TR_HAS_TRIES_CLASS = 'has-behaviour';
const TR_ACTIVE_CLASS = 'active';

combat.renderTable = function($target, logs) {
	const {createTag, showTries} = this;
	const $tries = createTag('div', {class: 'col-6'});

	var items = Object.keys(combatLogs).reduce((agg, key) => {
		return agg.concat([
			Object.assign({}, combatLogs[key], {key})
		])
	}, [])
		.sort((a, b) => {
			if (a.title > b.title) {
				return 1;
			}

			if (a.title < b.title) {
				return -1;
			}

			return 0;
		});

	const $tbody = createTag('tbody', {
		children: items.map((item, index) => {
			const {status, title, tries, lastSuccessfulRun} = item;

			const hasBehaviour = tries && tries.length > 0;
			// var triesSection= [];
			// if (combat.config.silentTries){
			// 	triesSection.push(
			// 		createTag('td', {class:'test_tries_span', children:
            //                 createTag('span', {children: tries===null? 'pending' : tries.length})}),
            //         createTag('td', {class:'test_name_span', children:
            //                 createTag('span', {children: title})}))
			// }else {
            //     triesSection.push(
            //     	createTag('td', {class:'test_tries_span'}),
            //         createTag('td', {class:'test_name_span', children:
            //                 createTag('span', {children: title})}))
			// }
            var $tr;
            if (window.silentTries) {
                 $tr = createTag('tr', {
                    class: hasBehaviour ? TR_HAS_TRIES_CLASS : '',
                    children: [
                        createTag('th', {scope: 'row', children: 1 + index}),
                        createTag('td', {
                            children:
                                createTag('div', {class: `icon icon--${status}`})
                        }),
                        createTag('td', {
                            class: 'test_name_span', children:
                                createTag('span', {children: title})
                        })
                    ]
                });
            }
            else {
                 $tr = createTag('tr', {
                    class: hasBehaviour ? TR_HAS_TRIES_CLASS : '',
                    children: [
                        createTag('th', {scope: 'row', children: 1 + index}),
                        createTag('td', {
                            children:
                                createTag('div', {class: `icon icon--${status}`})
                        }),
                        createTag('td', {
                            class: 'test_tries_span', children:
                                createTag('span', {children: tries === null ? 'pending' : tries.length})
                        }),
                        createTag('td', {
                            class: 'test_name_span', children:
                                createTag('span', {children: title})
                        })
                    ]
                });
            }

			if (hasBehaviour) {
				$tr.addEventListener('click', ({target}) => {
					const $prev = $tbody.querySelector(`.${TR_ACTIVE_CLASS}`);
					const parsedClassName = ` ${TR_ACTIVE_CLASS}`;

					if ($prev) {
						$prev.className = $prev.className.replace(parsedClassName, '')
					}

					target.closest('tr').className += parsedClassName;
					currentSliderIndex =null;
					showTries($tries, tries, lastSuccessfulRun, status);
				}, false);
			}

			return $tr;
		})
	});

	if(window.silentTries) {
        $target.append(
            createTag('div', {
                class: 'container', children:
                    createTag('div', {
                        class: 'row', children: [
                            createTag('div', {
                                class: 'col-6', children: [
                                    // createTag('button', { id: 'disable_slack', children : 'Disable slack notification (8 hours)'}),
                                    createTag('table', {
                                        class: 'table table-hover', children: [
                                            createTag('thead', {
                                                class: 'thead-light', children:
                                                    createTag('tr', {
                                                        class: 'thead-light', children: [
                                                            createTag('th', {css: 'col', scope: 'col', children: '#'}),
                                                            createTag('th', {
                                                                css: 'col',
                                                                scope: 'col',
                                                                children: 'Status'
                                                            }),
                                                            createTag('th', {
                                                                cas: 'col-10',
                                                                scope: 'col',
                                                                children: 'Test name'
                                                            })
                                                        ]
                                                    })
                                            }),
                                            $tbody
                                        ]
                                    })]
                            }),
                            $tries
                        ]
                    })
            })
        );
    }
    else {
        $target.append(
            createTag('div', {
                class: 'container', children:
                    createTag('div', {
                        class: 'row', children: [
                            createTag('div', {
                                class: 'col-6', children: [
                                    createTag('table', {
                                        class: 'table table-hover', children: [
                                            createTag('thead', {
                                                class: 'thead-light', children:
                                                    createTag('tr', {
                                                        class: 'thead-light', children: [
                                                            createTag('th', {css: 'col', scope: 'col', children: '#'}),
                                                            createTag('th', {
                                                                css: 'col',
                                                                scope: 'col',
                                                                children: 'Status'
                                                            }),
                                                            createTag('th', {
                                                                css: 'col',
                                                                scope: 'col',
                                                                children: 'Tries'
                                                            }),
                                                            createTag('th', {
                                                                cas: 'col-10',
                                                                scope: 'col',
                                                                children: 'Test name'
                                                            })
                                                        ]
                                                    })
                                            }),
                                            $tbody
                                        ]
                                    })]
                            }),
                            $tries
                        ]
                    })
            })
        );
	}
}

