package io.signare.openapigenerators.signare;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import org.openapitools.codegen.*;
import io.swagger.models.properties.*;

import com.google.common.collect.ImmutableMap;
import com.samskivert.mustache.Mustache;

import java.io.IOException;
import java.util.*;
import java.io.File;

import io.signare.openapigenerators.signare.templating.mustache.SignareCamelCaseLambda;

import org.openapitools.codegen.languages.AbstractGoCodegen;
import org.openapitools.codegen.model.ModelMap;
import org.openapitools.codegen.model.ModelsMap;
import org.openapitools.codegen.model.OperationsMap;
import org.openapitools.codegen.utils.StringUtils;


public class SignareHttpInfraGenerator extends AbstractGoCodegen implements CodegenConfig {

    // source folder where to write the files
    protected String sourceFolder = "";
    protected String apiVersion = "1.0.0";

    private static final String ACTIONS_FILE_PATH_PROPERTY_NAME = "generatedActionsFilePath";

    private List<String> actions = new ArrayList<>();

    /**
     * Configures the type of generator.
     *
     * @return the CodegenType for this generator
     * @see org.openapitools.codegen.CodegenType
     */
    public CodegenType getTag() {
        return CodegenType.OTHER;
    }

    /**
     * Configures a friendly name for the generator.  This will be used by the generator
     * to select the library with the -g flag.
     *
     * @return the friendly name for the generator
     */
    public String getName() {
        return "signare-http-infra";
    }


    /**
     * Provides an opportunity to inspect and modify operation data before the code is generated.
     */
    @SuppressWarnings("unchecked")
    @Override
    public OperationsMap postProcessOperationsWithModels(OperationsMap objs, List<ModelMap> allModels) {
        OperationsMap results = super.postProcessOperationsWithModels(objs, allModels);
        Map<String, Object> ops = (Map<String, Object>) results.get("operations");
        ArrayList<CodegenOperation> opList = (ArrayList<CodegenOperation>) ops.get("operation");
        // iterate over the operation and perhaps modify something
        for (CodegenOperation co : opList) {
            co.vendorExtensions.put("signare-action", co.operationIdOriginal);
            for (CodegenResponse r : co.responses) {
                // A valid response that is a string cannot be validated with the generated code. This only happens at the moment for an endpoint returning a file.
                if (r.is2xx) {
                    co.vendorExtensions.put("signare-response-type", r.dataType);
                    if (r.isString) {
                        co.vendorExtensions.put("signare-skip-code-generation", true);
                    }
                }
            }

            for (CodegenParameter param : co.allParams) {
                // Check if there is any String or an Array add a flag property in that case that will be used to modify handlers test generation behaviour
                if (param.isString || param.isArray) {
                    co.vendorExtensions.put("hasRequestValidation", true);
                }

                if (param.dataFormat == null) {
                    param.dataFormat = "SignareFreeText";
                }
                param.dataFormat = normalizeDataFormat(param.dataFormat);
            }
        }

        if (generateActions()) {
            generateActionsFile((String) this.additionalProperties().get(ACTIONS_FILE_PATH_PROPERTY_NAME), opList);
        }

        return results;
    }

    private boolean generateActions() {
        return this.additionalProperties().containsKey(ACTIONS_FILE_PATH_PROPERTY_NAME) && !this.additionalProperties().get(ACTIONS_FILE_PATH_PROPERTY_NAME).equals("");
    }

    @Override
    public ModelsMap postProcessModels(ModelsMap objs) {
        ModelsMap results = super.postProcessModels(objs);
        List<Map<String, Object>> modelList = (List) results.get("models");
        for (Map<String, Object> modelItem : modelList) {
            CodegenModel codeGenModel = (CodegenModel) modelItem.get("model");
            for (CodegenProperty prop : codeGenModel.vars) {
                if (prop.dataFormat != null) {
                    prop.dataFormat = normalizeDataFormat(prop.dataFormat);
                }
            }
        }

        return results;
    }

    public String normalizeDataFormat(String param) {
        return StringUtils.camelize(param);
    }

    /**
     * Returns human-friendly help for the generator.  Provide the consumer with help
     * tips, parameters here
     *
     * @return A string value for the help message
     */
    public String getHelp() {
        return "Generates a signare-http-infra client library.";
    }

    public SignareHttpInfraGenerator() {
        super();

        // set the output folder here
        outputFolder = "generated-code/signare-http-infra";

        /**
         * Models.  You can write model files using the modelTemplateFiles map.
         * if you want to create one template for file, you can do so here.
         * for multiple files for model, just put another entry in the `modelTemplateFiles` with
         * a different extension
         */
        modelTemplateFiles.put(
                "components_types.mustache", // the template to use
                "_types_generated.go");       // the extension for each file to write

        /**
         * Api classes.  You can write classes for each Api file with the apiTemplateFiles map.
         * as with models, add multiple entries with different extensions for multiple files per
         * class
         */
        apiTemplateFiles.put(
                "paths_handlers.mustache",
                "_http_handler_generated.go");

        /*
        apiTemplateFiles.put(
                "paths_handlers_test.mustache",
                "_http_handler_generated_test.go");
         */

        apiTemplateFiles.put(
                "paths_publisher.mustache",
                "_publisher_generated.go"
        );

        apiTemplateFiles.put(
                "paths_publisher_test.mustache",
                "_publisher_generated_test.go"
        );

        apiTemplateFiles.put(
                "paths_api_types.mustache",
                "_types_generated.go");

        /**
         * Template Location.  This is the location which templates will be read from.  The generator
         * will use the resource stream to attempt to read the templates.
         */
        templateDir = "signare-http-infra";

        /**
         * Api Package.  Optional, if needed, this can be used in templates
         */
        apiPackage = "org.openapitools.api";

        /**
         * Model Package.  Optional, if needed, this can be used in templates
         */
        modelPackage = "org.openapitools.model";

        /**
         * Reserved words.  Override this with reserved words specific to your language
         */
        reservedWords = new HashSet<String>(
                Arrays.asList(
                        "sample1",  // replace with static values
                        "sample2")
        );

        /**
         * Additional Properties.  These values can be passed to the templates and
         * are available in models, apis, and supporting files
         */
        additionalProperties.put("apiVersion", apiVersion);

        /**
         * Supporting Files.  You can write single files for the generator with the
         * entire object tree available.  If the input file has a suffix of `.mustache
         * it will be processed by the template engine.  Otherwise, it will be copied
         */

        supportingFiles.add(new SupportingFile("supporting_converters.mustache",
                "httpinfra",
                "signare_converters_generated.go")
        );

        supportingFiles.add(new SupportingFile("supporting_typestest.mustache",
                "httpinfra",
                "types_test.go")
        );

        /**
         * Language Specific Primitives.  These types will not trigger imports by
         * the client generator
         */
        languageSpecificPrimitives = new HashSet<String>(
                Arrays.asList(
                        "Type1",      // replace these with your types
                        "Type2")
        );
    }


    /**
     * Escapes a reserved word as defined in the `reservedWords` array. Handle escaping
     * those terms here.  This logic is only called if a variable matches the reserved words
     *
     * @return the escaped term
     */
    @Override
    public String escapeReservedWord(String name) {
        return "_" + name;  // add an underscore to the name
    }

    /**
     * Location to write model files.  You can use the modelPackage() as defined when the class is
     * instantiated
     */
    public String modelFileFolder() {
        return outputFolder + "/" + sourceFolder + "/" + modelPackage().replace('.', File.separatorChar);
    }

    /**
     * Location to write api files.  You can use the apiPackage() as defined when the class is
     * instantiated
     */
    @Override
    public String apiFileFolder() {
        return outputFolder + "/" + sourceFolder + "/" + apiPackage().replace('.', File.separatorChar);
    }

    /**
     * override with any special text escaping logic to handle unsafe
     * characters so as to avoid code injection
     *
     * @param input String to be cleaned up
     * @return string with unsafe characters removed or escaped
     */
    @Override
    public String escapeUnsafeCharacters(String input) {
        return input;
    }

    /**
     * Escape single and/or double quote to avoid code injection
     *
     * @param input String to be cleaned up
     * @return string with quotation mark removed or escaped
     */
    public String escapeQuotationMark(String input) {
        return input.replace("\"", "\\\"");
    }

    @Override
    public String toModelFilename(String name) {
        String fileName = super.toModelFilename(name);
        return "components_" + fileName.replaceFirst("model_", "");
    }

    @Override
    public String toApiFilename(String name) {
        String fileName = super.toApiFilename(name);
        return "api_" + fileName.replaceFirst("api_", "");
    }

    @Override
    protected ImmutableMap.Builder<String, Mustache.Lambda> addMustacheLambdas() {
        return super.addMustacheLambdas()
                .put("signareCamelCase", new SignareCamelCaseLambda());
    }

    private void generateActionsFile(String actionsOutputFilePath, ArrayList<CodegenOperation> operations) {
        // Get the actions list (they map 1 to 1 with operation IDs)
        for (CodegenOperation operation: operations) {
            actions.add(operation.operationIdOriginal);
        }
        // Create new YAML file
        Actions generatedActions = new Actions(actions);

        // Output actions to file
        ObjectMapper om = new ObjectMapper(new YAMLFactory());
        try {
            om.writeValue(new File(actionsOutputFilePath), generatedActions);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}

