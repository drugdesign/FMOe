import basis_set_exchange
import json
import sys


def number_of_functions(shells):
    cont_map = {}
    for sh in shells:
        ngeneral = len(sh['shell_coefficients'])

        is_spdf = len(sh['shell_angular_momentum']) > 1

        for am in sh['shell_angular_momentum']:
            ncont = ngeneral if not is_spdf else 1

            if am not in cont_map:
                cont_map[am] = ncont
            else:
                cont_map[am] = cont_map[am] + ncont

    s = 0
    for am, ncont in cont_map.items():
        s += ((am + 1) * (am + 2) // 2) * ncont

    return s


def main():
    d = {}
    for basisset_name, meta in basis_set_exchange.get_metadata().items():
        basisset = {}
        elements = basis_set_exchange.get_basis(basisset_name)['basis_set_elements']

        try:
            for Z in meta['versions'][meta['latest_version']]['elements']:
                element = basis_set_exchange.lut.element_data_from_Z(Z)[0].capitalize()
                shells = elements[Z].get('element_electron_shells')
                if shells is None:
                    raise ValueError('no shells: {}/{}'.format(basisset_name, element))
                basisset[element] = number_of_functions(shells)

        except ValueError as e:
            print(e, file=sys.stderr)
        else:
            d[basisset_name] = basisset

    json.dump(d, sys.stdout)


if __name__ == "__main__":
    main()
